package db

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/tiantaozhang/go-blog/util"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strings"
	"time"
)

const (
	salt = "1234"
)

func Pwd(pass string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := fmt.Sprintf("%v", r.Intn(8999)+1000)
	m := md5.New()
	m.Write([]byte(pass))
	md5P := m.Sum(nil)
	bs64 := base64.StdEncoding.EncodeToString( []byte(string(md5P) + salt))
	beelog.Debug("pwd erypto:%v", bs64)
	return bs64
}

func ComparePwd(passTC, passC string) error {
	bs64,err := base64.StdEncoding.DecodeString(passC)
	if err !=nil {
		beelog.Error("base64 decodestring(%v)-->err(%v)",passTC,err)
		return err
	}
	salt := []byte(bs64)[len(bs64)-4:]
	m := md5.New()
	m.Write([]byte(passTC))
	md5P := m.Sum(nil)
	cptP := string(md5P) + string(salt)
	if strings.EqualFold(cptP, string(bs64)) {
		return nil
	}
	return fmt.Errorf("%v", "pwd error")
}

func AddUsers(users []*User) ([]map[string]interface{}, error) {
	m := []map[string]interface{}{}
	for _, u := range users {
		intid, id, err := NewUid()
		if err != nil {
			return nil, err
		}
		u.Id = id
		//u.Pwd = Pwd(u.Pwd)
		u.Pwd=Encrypt(u.Pwd,"")
		u.Status = US_NORMAL
		u.Time = util.TimeM()
		u.Last = util.TimeM()
		defaultusr:=fmt.Sprintf("%v",intid+1234567)
		if u.User==nil {
			u.User=append(u.User,defaultusr)
		}else{
			tmpUsr:=u.User
			u.User=[]string{defaultusr}
			u.User=append(u.User,tmpUsr...)
		}

		update := bson.M{}
		update["$setOnInsert"] = u
		_,err = C(CN_USER).Upsert(bson.M{"usrname": bson.M{"$in": u.User}, "pwd": u.Pwd, "type": u.Type, "status": u.Status}, update)
		if err == nil {
			m = append(m, map[string]interface{}{
				"id":      u.Id,
				"usrname": u.User,
				"type":    u.Type,
				"time":    u.Time,
			})
		} else {
			return nil,err
		}

	}
	return m,nil
}
