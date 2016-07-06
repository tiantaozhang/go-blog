package db

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/tiantaozhang/go-blog/util"
	"github.com/julienschmidt/httprouter"
)

var user1, user2, user3 *User

func init() {
	RemoveUser()
	user1=new(User)
	user2=new(User)
	user3=new(User)
	user1.Pwd = "123456"
	user2.Pwd = "123456"
	user3.Pwd = "123456"
	_, err := AddUsers([]*User{user1, user2, user3})
	if err != nil {
		panic(err)
	}
	//beelog.Debug("AddUsers resm:(%v)",util.S2Json(resm) )
	beelog.Debug("user1:(%v)",util.S2Json(user1))
	beelog.Debug("user2:(%v)",util.S2Json(user2))
	beelog.Debug("user3:(%v)",util.S2Json(user3))
}

func RemoveUser() error{
	if _,err:=C(CN_USER).RemoveAll(nil);err!=nil{
		beelog.Error("user removeall err:%v",err)
		return err
	}
	return nil
}

func TestLoginAndOut(t *testing.T) {

	router := httprouter.New()
	router.POST("/login", Login)
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
	// UsrInfo.User = "u1"
	// UsrInfo.Pwd="123456"
	// ju, err = json.Marshal(UsrInfo)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// beelog.Debug("ju:(%v)", string(ju))
	r2, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"user":"u1","pwd":"123456"}`))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)
	if w2.Code != http.StatusNotFound {
		t.Error("w2.code:", w2.Code)
		return
	}

}
