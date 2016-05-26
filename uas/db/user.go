package db

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/tiantaozhang/go-blog/util"
	"gopkg.in/mgo.v2"
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
	bs64 := base64.StdEncoding.EncodeToString(md5P + salt)
	fmt.Printf("pwd erypto:%v", bs64)
	return bs64
}

func ComparePwd(passTC, passC string) error {
	bs64 := base64.StdEncoding.DecodeString(passC)
	salt := []byte(bs64)[len(bs64)-4:]
	m := md5.New()
	m.Write(passTC)
	md5P := m.Sum(nil)
	cptP := md5P + salt
	if strings.EqualFold(cptP, string(bs64)) {
		return nil
	}
	return fmt.Errorf("%v", "pwd error")
}

func AddUsers(users []*User) ([]map[string]interface{}, error) {
	m := []map[string]interface{}{}
	for _, u := range users {
		_, id, err := NewUid()
		if err != nil {
			return nil, err
		}
		u.Id = id
		u.Pwd = Pwd(u.Pwd)
		u.Status = US_NORMAL
		u.Time = util.TimeM()
		u.Last = util.TimeM()

		update := bson.M{}
		update["$setOnInsert"] = u
		err = C(CN_USER).Upsert(bson.M{"usrname": {"$in": u.UsrName}, "pwd": u.Pwd, "type": u.Type, "status": u.Status}, update)
		if err == nil {
			m = append(m, map[string]interface{}{
				"id":      u.Id,
				"usrname": u.UsrName,
				"type":    u.Type,
				"time":    u.Time,
			})
		} else {
			return err
		}

	}

}
