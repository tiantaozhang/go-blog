package db

import (
	"gopkg.in/mgo.v2"
)

func init() {
	defer func() {
		recover()
	}()

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	C = session.DB("blog")
}
