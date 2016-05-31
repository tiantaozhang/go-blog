package session

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/tiantaozhang/go-blog/util"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"
)

const (
	EXPIRES = 10
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

func (m *ManagerSession) NewSession(key string, values map[interface{}]interface{}) *Session {
	return &Session{
		Sid:     m.GenSid(),
		Key:     key,
		Values:  values,
		Expires: m.Expires,
	}
}

func (m *ManagerSession) Start(w http.ResponseWriter, r *http.Request) (session Session) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	cookie, err := r.Cookie(m.Name)
	if err != nil || cookie == "" {
		session = m.NewSession(m.Name, nil)
		m.Set(session)
		cookie = http.Cookie{Name: m.Name, Value: url.QueryEscape(s.Sid), Path: "/", HttpOnly: true, MaxAge: m.Expires}
		http.SetCookie(w, cookie)
		return
	} else {
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			log.Printf("in start,QueryUnescape sid:(%v),err:(%v)", sid, err)
			return
		}
		if session, err = m.Get(sid); err != nil {
			log.Printf("in start,get session by sid:(%v),err:(%v)", sid, err)
			return
		}
		return
	}

}

func (m *ManagerSession) Get(sid string) (Session, error) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	if m.SM[sid].Value == nil {
		return nil, fmt.Errorf("can't find session by sid:%v", sid)
	}
	if reflect.TypeOf(m.SM[sid]).Kind() != reflect.Struct {
		return nil, fmt.Errorf("type's kind is not struct,but %v", reflect.TypeOf(m.SM[sid]).Kind())
	}
	m.SM[sid].Value.(*Session).Expires = m.Expires
	s := m.SM[sid].Value.(*Session)
	//garbage recover
	for e := m.SL.Front(); e != nil; e = e.Next() {
		if sid == e.Value.(*Session).Sid {
			e.Value = s
			m.SL.MoveToBack(e)
			break
		}
	}

	return *s, nil
}

func (m *ManagerSession) Listen() {

}

func (m *ManagerSession) GetByKey() {

}

func (m *ManagerSession) Set(s *Session) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if s.Sid == "" || s.Key == m.Name {
		log.Printf("session err:(%v)\n", util.S2Json(s))
		return fmt.Errorf("session err:(%v)", util.S2Json(s))
	}
	if m.SM[s.Sid].Value == nil {
		//not exist
		m.SL.PushBack(s)
	} else {
		for e := m.SL.Front(); e != nil; e = e.Next() {
			if s.Sid == e.Value.(*Session).Sid {
				e.Value.(*Session).Expires = m.Expires
				m.SL.MoveToBack(e)
				break
			}
		}
	}

	m.SM[s.Sid] = &list.Element{Value: s}
	return nil

}

func (m *ManagerSession) Del() {

}

func (m *ManagerSession) GenSid() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
