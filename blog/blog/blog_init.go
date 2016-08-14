package blog

import (
	"time"

	"github.com/tiantaozhang/go-blog/blog/db"
	"github.com/tiantaozhang/go-blog/logs"
	"gopkg.in/mgo.v2"
)

func init() {
	defer func() {
		recover()
	}()
	count := 0
	for {
		session, err := mgo.Dial("cny:123@127.0.0.1:27017")
		if err != nil {
			panic(err)
		}
		if err = session.Ping(); err != nil {
			logs.Beelog.Error("session ping err:%v", err)
			session.Close()
		} else {
			db.C = session.DB("blog").C
			break
		}
		time.Sleep(time.Second)
		count++
		if count > 5 {
			panic(err)
		}
	}
	logs.Beelog.Info("blog db init success")

}
