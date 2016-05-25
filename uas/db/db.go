package db

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const (
	CN_SEQUENCE  = "sequence"      //the sequence collection name
	CN_USER      = "user"          //the user collection name.
	CN_ANONYMOUS = "uas_anonymous" //the anonymous user record
)

var C = func(name string) *mgo.Collection {
	panic("the user db collection is not initial")
}
