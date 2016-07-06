package db

type User struct {
	Id     string                            `bson:"_id" json:"_id" `
	User   []string                          `bson:"user" json:"user,omitempty"`     //the user name list
	Pwd    string                            `bson:"pwd" json:"pwd,omitempty"`       //the user password
	Attrs  map[string]map[string]interface{} `bson:"attrs" json:"attrs,omitempty"`   //the user attribute list split by group
	Type   int                               `bson:"type" json:"type,omitempty"`     //type
	Status int                               `bson:"status" json:"status,omitempty"` //the user status
	Last   int64                             `bson:"last" json:"last,omitempty"`     //the last updated time
	Time   int64                             `bson:"time" json:"time,omitempty"`     //the create time

}

const (
	UT_ADMIN   = 10 //the admin
	UT_BLOGGER = 20 //the blog owner
)

const (
	US_NORMAL = 10
	US_DELETE = -1
)

type Sequence struct {
	Id  string `bson:"_id" json:"id"`  //the sequenc id.
	Val uint64 `bson:"val" json:"val"` //the current sequene value
}
