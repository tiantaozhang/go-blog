package db

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
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
