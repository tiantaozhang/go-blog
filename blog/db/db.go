package db

import (
	"fmt"
	"time"

	"github.com/tiantaozhang/go-blog/logs"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CN_SEQUENCE = "sequence" //the sequence collection name
	CN_BLOG     = "blog"     //the user collection name.
)

var C = func(name string) *mgo.Collection {
	panic("the blog db collection is not initial")
}

const (
	SQN_BLOG = "blog" //the sequence id
)

type Sequence struct {
	Id  string `bson:"_id" json:"id"`  //the sequenc id.
	Val uint64 `bson:"val" json:"val"` //the current sequene value
}

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

func InitDb() {
	for {
		session, err := mgo.Dial("cny:123@127.0.0.1:27017")
		if err != nil {
			logs.Beelog.Error("mgo dial (%v) err(%v)", "cny:123/127.0.0.1:27017", err)
			panic(err)
		}
		if err = session.Ping(); err != nil {
			logs.Beelog.Error("session ping err:%v", err)
			session.Close()
			time.Sleep(time.Second)
			continue
		}
		C = session.DB(CN_BLOG).C
		break

	}

}
