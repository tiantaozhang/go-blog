package session

import (
	"container/list"
	"time"

	//"github.com/tiantaozhang/go-blog/util"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSession(t *testing.T) {
	Ms := &ManagerSession{
		SL:      list.New(),
		SM:      make(map[interface{}]*list.Element),
		Name:    "gosid",
		Expires: 2,
	}

	var err error
	var sid string
	//var sess Session
	//ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := Ms.Start(w, r)
	if err != nil {
		t.Error(err.Error())
		return
	}
	sid = sess.Sid

	//test session again
	cookie := http.Cookie{Name: Ms.Name, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: Ms.Expires}
	r.AddCookie(&cookie)
	sess1, err := Ms.Start(w, r)
	if err != nil {
		t.Error(err.Error())
		return
	} else if sess1.Sid != sid {
		t.Error("not equal")
	}

	if sess.Key != Ms.Name {
		t.Error("init error")
		return
	} else {
		FmtPrintf("init session,sess.Key:%v,sess.sid:%v", sess.Key, sess.Sid)
	}

	FmtPrintf("Ms.SM:%v,Ms.SM[sess.Sid]:%v", Ms.SM, Ms.SM[sess.Sid])

	if sess.Sid != Ms.SM[sess.Sid].Value.(*Session).Sid {
		t.Error("err")
		return
	}

	Ms.Listen()
	if sess, err = Ms.Get(sid); err != nil {
		t.Error("get session err:%v", err.Error())
	} else {
		FmtPrintf("session id:%v,key:%v", sess.Sid, sess.Key)
	}
	//test Set equal key
	news := Ms.NewSession(nil)
	Ms.Set(news)

	if sess, err = Ms.Get(news.Sid); err != nil {
		t.Error("get new %s", err.Error())
	} else {
		FmtPrintf("key:%s,sid:%s", sess.Key, sess.Sid)
	}

	//the first sid

	//r.Header.Set("Set-Cookie", )
	r1, _ := http.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()
	r1.AddCookie(&cookie)
	sessrtn, err := Ms.Start(w1, r1)
	if sessrtn.Key == "" {
		t.Error("sessrtn Key nil")
	} else if sessrtn.Key != Ms.Name || sessrtn.Sid != sid {
		t.Error("key:%s", sessrtn.Key)
	} else {
		fmt.Println("key and sid marry")
	}

	if err := Ms.Del(sid); err != nil {
		t.Error("del sid:(%v) error:(%v)", sess.Sid, err.Error())
	} else {
		fmt.Println("del sess.sid:(%v) right", sess.Sid)
	}
	if err := Ms.Del("123"); err == nil {
		t.Error("del 123 err")
	}
	//del and request

	sessnil, err := Ms.Start(w, r)
	if sessnil.Sid == sid {
		t.Error("new sess sid:(%v) is same as old sid:(%v)", sess.Key, sid)
	} else {
		FmtPrintf("new sessnil sid:%v", sessnil.Sid)
	}

	sessionNew, err := Ms.Start(w1, r1)
	if err != nil {
		t.Error(err.Error())
	}
	snew, err := Ms.Get(sessionNew.Sid)
	if snew.Sid == "" || snew.Sid == sid {
		t.Error("new sid:", snew.Sid)
		return
	}

	time.Sleep(7 * time.Second)

	if s, err := Ms.Get(sessionNew.Sid); err == nil {
		t.Error("err,sid:", s.Sid)
		return
	}

}
