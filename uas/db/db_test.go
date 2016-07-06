package db

import (
	"gopkg.in/mgo.v2"
	"time"
)

func init() {
	defer func() {
		recover()
	}()
	count:=0
for{
	session, err := mgo.Dial("cny:123@127.0.0.1:27017/blog")
	if err != nil {
		panic(err)
	}
	if err=session.Ping();err!=nil{
		beelog.Error("session ping err:%v",err)
	}else{
		C = session.DB("blog").C
		break
	}
	time.Sleep(time.Second)
	count++
	if count>5 {
		panic(err)
	}
}

}
