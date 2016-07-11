package db

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tiantaozhang/go-blog/uas/session"
	"github.com/tiantaozhang/go-blog/util"
	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego/logs"

	"github.com/julienschmidt/httprouter"
)

var beelog *logs.BeeLogger

func init() {

	beelog = logs.NewLogger(0)
	beelog.SetLogger("console", "")

	beelog.EnableFuncCallDepth(true)
	beelog.SetLevel(logs.LevelDebug)
	beelog.SetLogFuncCallDepth(2)

}

func CreateUsers(u []*User) ([]map[string]interface{}, error) {

	if err := checkUser(u); err != nil {
		return nil, err
	}

	return AddUsers(u)
}

func checkUser(user []*User) error {
	if user == nil {
		return errors.New("u is nil")
	}
	for _, u := range user {
		if u.User == nil {
			return fmt.Errorf("%v", "please input your usrName")
		}
		if u.Pwd == "" {
			return fmt.Errorf("%v", "please input your pwd")
		}
		if u.Type != UT_ADMIN && u.Type != UT_BLOGGER {
			//default blogger
			u.Type = UT_BLOGGER
		}
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beelog.Debug("---------Login-------- ")

	if c, err := r.Cookie(session.SM.Name);err==nil{
		sid, _ := url.QueryUnescape(c.Value)
		//sid:=c.Value
		if _,err=session.SM.Get(sid);err!=nil{
			beelog.Error("login cookie not nil,sid(%v)-->err(%v)",sid,err)
			//http.Error(w,err.Error(),http.StatusNotFound)
			//return
		}else{
			beelog.Debug("cookie not nil, c.value:%v",sid)
			http.Redirect(w, r, "http://www.baidu.com", http.StatusOK)
			return
		}

	}

	if err := r.ParseForm(); err != nil {
		beelog.Error("login parseform err(%v)", err)
		return
	}
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		beelog.Error("login readall err:(%v)", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	beelog.Debug("request body:(%v)", string(res))
	var UsrInfo struct {
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	}
	if err := json.Unmarshal(res, &UsrInfo); err != nil {
		beelog.Error("login Unmarshal err(%v)", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user := UsrInfo.User
	pwd := UsrInfo.Pwd
	//user := r.Form["user"]
	// user := ps.ByName("user")
	// pwd := ps.ByName("pwd")
	beelog.Debug("user(%v) pwd(%v) ps(%v)", UsrInfo.User, UsrInfo.Pwd, util.S2Json(ps))

	if user == "" || pwd == "" {
		beelog.Error("username(%v) or pwd(%v) nil", user, pwd)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var u User
	if err := C(CN_USER).Find(bson.M{"user": user, "pwd": Encrypt(pwd, "")}).One(&u); err != nil {
		beelog.Error("query(%v) err(%v)", util.S2Json(bson.M{"user": user, "pwd": Encrypt(pwd, "")}), err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	//?
	r.Form.Set("uid", u.Id)
	if sess, err := do_login(w, r); err != nil {
		beelog.Error("do_login --> err(%v)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// redirect
		beelog.Debug("login user(%v) success session sid(%v) value(%v) key(%v)", user, sess.Sid,sess.Values,sess.Key)
		http.Redirect(w, r, "http://www.baidu.com", http.StatusOK)
	}

	return
}

func Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beelog.Debug("------------logout----------")

	r.ParseForm()

	c, err := r.Cookie(session.SM.Name)
	if err != nil {
		beelog.Error("r.cookie (%v)-->err(%v)", session.SM.Name, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	sid, _ := url.QueryUnescape(c.Value)
	err = session.SM.Del(sid)
	if err != nil {
		beelog.Error("Del sid(%v) err(%v)", sid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		beelog.Error("login readall err:(%v)", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	beelog.Debug("request body:(%v)", string(res))
	var UsrInfo struct {
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	}
	if err := json.Unmarshal(res, &UsrInfo); err != nil {
		beelog.Error("login Unmarshal err(%v)", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	beelog.Debug("user(%v) logout",UsrInfo.User )
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}

func do_login(w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	return session.SM.Start(w, r)
}

//4 bits
func Encrypt(pwd string, salt string) string {
	salt = "1234"
	h := md5.New()
	h.Write([]byte(pwd))
	md5pwd := h.Sum(nil)

	base64pwd := base64.StdEncoding.EncodeToString([]byte(string(md5pwd) + salt))

	return base64pwd
}

func CheckCrypt(base64pwd, pwd string) bool {
	saltpwd, err := base64.StdEncoding.DecodeString(base64pwd)
	if err != nil {
		beelog.Error("base64 decodestring(%v)-->err(%v)", base64pwd, err)
		return false
	}
	md5pwd := string([]byte(saltpwd)[:len(saltpwd)-4])
	h := md5.New()
	h.Write([]byte(pwd))
	eptpwd := h.Sum(nil)
	if string(eptpwd) == md5pwd {
		return true
	}
	return false
}
