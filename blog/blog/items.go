package blog

type Blog struct {
	Id      string                 `bson:"_id" json:"_id"`
	Title   string                 `bson:"title" json:"title"`
	Content string                 `bson:"content" json:"content"`
	Author  string                 `bson:"author" json:"author"`
	Replies int                    `bson:"reply" json:"reply"`
	Views   int                    `bson:"views" json:"views"`
	Tags    []string               `bson:"tags" json:"tags"`
	Ext     map[string]interface{} `bson:"ext" json:"ext"`
	Type    int                    `bson:"type" json:"type"`
	Status  int                    `bson:"status" json:"status"`
	Time    int64                  `bson:"time" json:"time"`
	Last    int64                  `bson:"last" json:"last"`
}

const (
	BT_ORIGIN      = 10
	BT_REPRODUCT   = 20
	BT_TRANSLATION = 30
)

const (
	BS_NORMAL = 10
	BS_DELETE = -1
)
