package session

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego/logs"
	"github.com/tiantaozhang/goColorChange"
)

var beelog *logs.BeeLogger

var SM *ManagerSession = &ManagerSession{
	SL:      list.New(),
	SM:      make(map[interface{}]*list.Element),
	Name:    "gosid",
	Expires: EXPIRES,
}

func init() {
	SM.Listen()

	beelog = logs.NewLogger(0)
	beelog.SetLogger("console", "")
	beelog.SetLevel(logs.LevelDebug)
	beelog.EnableFuncCallDepth(true)
}

const (
	EXPIRES = 10
)

const (
	GETERRNOTFOUND = "can't find session by sid"
)

type ManagerSession struct {
	Lock    sync.RWMutex
	SM      map[interface{}]*list.Element //storage sessions
	SL      *list.List                    //for gc
	Name    string
	Expires int
}

type Session struct {
	Lock    sync.RWMutex
	Sid     string                      //sid
	Key     string                      //correspond Manager Name
	Values  map[interface{}]interface{} //value
	Expires int                         //expire time
}

func (m *ManagerSession) NewSession(values map[interface{}]interface{}) *Session {
	if values == nil {
		values = make(map[interface{}]interface{})
	}

	return &Session{
		Sid:     m.GenSid(),
		Key:     m.Name,
		Values:  values,
		Expires: m.Expires,
	}
}

func (m *ManagerSession) Start(w http.ResponseWriter, r *http.Request) (session *Session, err error) {
	// m.Lock.Lock()
	// defer m.Lock.Unlock()
	cookie, err := r.Cookie(m.Name)
	if err != nil || cookie.Value == "" {
		beelog.Warn("session start if")
		cookie := new(http.Cookie)
		cookie, session, err = dealSessionHttp(m, r)
		http.SetCookie(w, cookie)
		return
	} else {
		beelog.Warn("session start else")
		var sid string
		sid, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			beelog.Debug("in start,QueryUnescape sid:(%v),err:(%v)", sid, err)
			return
		}
		if session, err = m.Get(sid); err != nil {
			beelog.Debug("in start,get session by sid:(%v),err:(%v)", sid, err)
			if err.Error() != GETERRNOTFOUND {
				return
			}
			cookie := new(http.Cookie)
			cookie, session, err = dealSessionHttp(m, r)
			if err != nil {
				beelog.Debug("in start dealSessionHttp,err:(%v)", err)
				return
			}
			http.SetCookie(w, cookie)

		}
		return
	}

}

func dealSessionHttp(m *ManagerSession, r *http.Request) (*http.Cookie, *Session, error) {
	session := m.NewSession(nil)
	var err error
	if err = r.ParseForm(); err != nil {
		beelog.Debug("manager start,parseForm,err:(%v)", err)
		return nil, nil, err
	}
	//uid & token=xxxx
	//user := r.FormValue("user")
	user := r.FormValue("user")
	session.Values["user"] = user
	//find uid by user and pwd
	//session.Values["uid"] = uid
	session.Values["token"] = GenToken()

	beelog.Debug("session:sid:(%v),values:(%v)", session.Sid, session.Values)
	if err = m.Set(session); err != nil {
		beelog.Debug("manager start,set session,err:(%v)", err)
		return nil, nil, err
	}
	beelog.Debug("session sid:(%v),session values:(%v)", session.Sid, session.Values)
	cookie := http.Cookie{Name: m.Name, Value: url.QueryEscape(session.Sid), Path: "/", HttpOnly: true, MaxAge: m.Expires}
	return &cookie, session, nil
}

func (m *ManagerSession) Get(sid string) (*Session, error) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	if m.SM[sid] == nil {
		//return nil, fmt.Errorf("can't find session by sid:%v", sid)
		return nil, fmt.Errorf("%v", GETERRNOTFOUND)
	}
	if reflect.TypeOf(*m.SM[sid]).Kind() != reflect.Struct {
		return nil, fmt.Errorf("type's kind is not struct,but %v", reflect.TypeOf(m.SM[sid]).Kind())
	}
	m.SM[sid].Value.(*Session).Expires = m.Expires
	s := m.SM[sid].Value.(*Session)
	//garbage recover
	m.SL.MoveToBack(m.SM[sid])

	return s, nil
}

func CallStatck() string {
	buf := make([]byte, 1024000)
	blen := runtime.Stack(buf, true)
	return string(buf[0:blen])
}

func (m *ManagerSession) Listen() {

	go func() {

		defer func() {
			if err := recover(); err != nil {
				beelog.Debug("the listen is panic->%v, the call stack is \n%v\n", err, CallStatck())
			}
		}()

		t := time.NewTicker(time.Duration(m.Expires) * time.Second)
		for {
			select {
			case <-t.C:
				m.Lock.RLock()
				var n *list.Element
				for e := m.SL.Front(); e != nil; e = n {
					n = e.Next()
					if e.Value.(*Session).Expires == 0 {
						FmtPrintf("listen gc-->delete sid:%v", e.Value.(*Session).Sid)
						m.Lock.RUnlock()
						m.Lock.Lock()
						m.SL.Remove(e)
						delete(m.SM, e.Value.(*Session).Sid)
						m.Lock.Unlock()
						m.Lock.RLock()
					} else {
						e.Value.(*Session).Expires--
						//m.SM[e.Value.(*Session).Sid].Value.(*Session).Expires--
						FmtPrintf("listen gc-->expires--:%v,sid:%v", e.Value.(*Session).Expires, e.Value.(*Session).Sid)
					}
				}
				m.Lock.RUnlock()
			}
		}
	}()

	//time.AfterFunc(time.Duration(m.Expires)*time.Second, func() { m.Listen() })
}

// func (m *ManagerSession) GetByKey(sid interface{}) {

// }

func (m *ManagerSession) Set(s *Session) error {

	m.Lock.Lock()
	defer m.Lock.Unlock()

	return m.set(s)

}

func (m *ManagerSession) set(s *Session) error {
	if s.Sid == "" || s.Key != m.Name {
		beelog.Error("session err:s.sid:(%v),s.key:%v\n", s.Sid, s.Key)
		return fmt.Errorf("session err:s.sid:(%v),s.key:(%v)", s.Sid, s.Key)
	}
	if m.SM[s.Sid] == nil {
		//not exist
		beelog.Warn("set function,s,m.SM[%v] is nil,pushback", s.Sid)
		m.SL.PushBack(s)
	} else {

		m.SL.MoveToBack(m.SM[s.Sid])
	}

	m.SM[s.Sid] = &list.Element{Value: s}
	return nil
}

func (m *ManagerSession) Del(sid string) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if m.SM[sid] == nil {
		return fmt.Errorf("not exist session's sid=%v", sid)
	}
	m.SL.Remove(m.SM[sid])
	delete(m.SM, sid)

	return nil
}

func (m *ManagerSession) GenSid() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GenToken() string {
	return bson.NewObjectId().Hex()
}

func FmtPrintf(format string, args ...interface{}) {

	goColorChange.ChangeColorAndStyle(goColorChange.Underline, goColorChange.Green, goColorChange.None)
	log.Printf(format, args...)
	log.Printf("\n")
	goColorChange.ResetColor()
}

func ChangeColorAndStyle(style goColorChange.Style, fg goColorChange.Color, bg goColorChange.Color) {
	goColorChange.ChangeColorAndStyle(style, fg, bg)
}

func ResetColor() {
	goColorChange.ResetColor()
}
