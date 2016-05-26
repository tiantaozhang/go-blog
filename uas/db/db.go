package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CN_SEQUENCE  = "sequence"      //the sequence collection name
	CN_USER      = "user"          //the user collection name.
	CN_ANONYMOUS = "uas_anonymous" //the anonymous user record
)

var C = func(name string) *mgo.Collection {
	panic("the user db collection is not initial")
}

const (
	SQN_USR   = "user" //the sequence id
	SQN_GROUP = "uas_group"
)

func QuerySequence() (uint64, error) {
	var seq Sequence
	_, err := C(CN_SEQUENCE).Find(bson.M{"_id": SQN_USR}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson{"val": 1}},
		Upsert:    true,
		ReturnNew: true,
	}, &seq)
	return seq.Val, err
}

func NewUid() (uint64, string, error) {
	uid, err := QuerySequence(SQN_USR)
	return uid, fmt.Sprintf("u%v", uid), err
}
