package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CN_SEQUENCE = "sequence" //the sequence collection name
	CN_BLOG     = "blog"     //the user collection name.
)

var C = func(name string) *mgo.Collection {
	panic("the user db collection is not initial")
}

const (
	SQN_BLOG = "blog" //the sequence id
)

func QuerySequence() (uint64, error) {
	var seq Sequence
	_, err := C(CN_SEQUENCE).Find(bson.M{"_id": SQN_BLOG}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"val": 1}},
		Upsert:    true,
		ReturnNew: true,
	}, &seq)
	return seq.Val, err
}

func NewBid() (uint64, string, error) {
	bid, err := QuerySequence()
	return bid, fmt.Sprintf("b%v", bid), err
}
