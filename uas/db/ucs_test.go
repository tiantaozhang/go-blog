package db

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/tiantaozhang/go-blog/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	//	"regexp"
	"github.com/tiantaozhang/go-blog/uas/session"
	"regexp"
	//	"net/url"
)

var user1, user2, user3 *User

func init() {
	RemoveUser()
	user1 = new(User)
	user2 = new(User)
	user3 = new(User)
	user1.Pwd = "123456"
	user2.Pwd = "123456"
	user3.Pwd = "123456"
	_, err := AddUsers([]*User{user1, user2, user3})
	if err != nil {
		panic(err)
	}
	//beelog.Debug("AddUsers resm:(%v)",util.S2Json(resm) )
	beelog.Debug("user1:(%v)", util.S2Json(user1))
	beelog.Debug("user2:(%v)", util.S2Json(user2))
	beelog.Debug("user3:(%v)", util.S2Json(user3))
}

func RemoveUser() error {
	if _, err := C(CN_USER).RemoveAll(nil); err != nil {
		beelog.Error("user removeall err:%v", err)
		return err
	}
	return nil
}

func TestLoginAndOut(t *testing.T) {

	router := httprouter.New()
	router.POST("/login", Login)
	router.POST("/logout",Logout)
	//test paras
	type UsrInfoStruct struct {
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	}
	var UsrInfo UsrInfoStruct
	UsrInfo.User = "u1"
	ju, err := json.Marshal(UsrInfo)
	if err != nil {
		t.Error(err)
		return
	}

	r1, _ := http.NewRequest("POST", "/login", strings.NewReader(string(ju)))
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, r1)
	if w1.Code != http.StatusFound {
		t.Error("w1.code:", w1.Code)
		return
	}
	//
	r2, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"user":"u1","pwd":"123456"}`))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)
	if w2.Code != http.StatusNotFound {
		t.Error("w2.code:", w2.Code)
		return
	}
	//test normal
	UsrInfo.User = user1.User[0]
	UsrInfo.Pwd = "123456"
	ju, err = json.Marshal(UsrInfo)
	r3, _ := http.NewRequest("POST", "/login", strings.NewReader(string(ju)))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Error("code not 200 err,code is %v", w3.Code)
		return
	}
	beelog.Debug("cookie:%v", w3.Header())
	cookie := w3.Header().Get("Set-Cookie")

	cookies := strings.Split(cookie, ";")
	beelog.Debug("header cookies:%v", cookies)
	rsid := regexp.MustCompile("gosid=")
	var sid string
	for _, c := range cookies {
		s := rsid.FindString(c)
		if s != "" {
			sids := strings.Split(c, "=")
			sid = string([]byte(c)[len(sids[0])+1:])
			beelog.Debug("gosid=(%v)", sid)
			break
		}
	}
	UsrInfo.User = user1.User[0]
	UsrInfo.Pwd = ""
	ju, err = json.Marshal(UsrInfo)
	r4, _ := http.NewRequest("POST", "/login", strings.NewReader(string(ju)))

	httpcookie := http.Cookie{Name: session.SM.Name, Value: sid, Path: "/", HttpOnly: true, MaxAge: session.SM.Expires}
	r4.AddCookie(&httpcookie)
	w4 := httptest.NewRecorder()

	router.ServeHTTP(w4, r4)
	if w4.Code != http.StatusOK {
		t.Error("code not 200 err,code is %v", w4.Code)
		return
	}

	r5,_:=http.NewRequest("POST","/logout", strings.NewReader(string(ju)))
	w5:=httptest.NewRecorder()

	router.ServeHTTP(w5,r5)
	if w5.Code!=http.StatusForbidden {
		t.Error("w5 err,",w5.Code)
		return
	}

	r5.AddCookie(&httpcookie)
	w6:=httptest.NewRecorder()
	router.ServeHTTP(w6,r5)
	if w6.Code!=http.StatusOK {
		t.Error("addcookie,w6 err,code:()",w6.Code)
		return
	}


}
