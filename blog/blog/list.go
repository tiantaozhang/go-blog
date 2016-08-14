package blog

import (
	"fmt"
	"strings"

	"github.com/tiantaozhang/go-blog/logs"
	"github.com/tiantaozhang/go-blog/util"

	"gopkg.in/mgo.v2/bson"

	"github.com/tiantaozhang/go-blog/blog/db"
)

//sort can be time,view
func ListBlog(uid, author string, page, pageCount int, bid, title string, tags []string, btype []int, sort string, fuzzyTitle bool, total bool) (map[string]interface{}, error) {

	if page == 0 || pageCount == 0 {
		logs.Beelog.Error("page(%v) pageCount(%v) is nil", page, pageCount)
		return nil, fmt.Errorf("page(%v) pageCount(%v) is nil", page, pageCount)
	}

	query := bson.M{}
	if uid != "" {
		query["uid"] = uid
	}
	if author != "" {
		query["author"] = author
	}
	if bid != "" {
		query["_id"] = bid
	}
	if title != "" {
		if fuzzyTitle != false {
			query["title"] = bson.M{"$regex": bson.M{"title": title}}
		} else {
			query["title"] = bson.M{"title": title}
		}
	}
	if btype != nil {
		query["type"] = bson.M{"$in": btype}
	}

	query["status"] = BS_NORMAL

	sortString := []string{}
	if sort == "" {
		sortString = []string{"-time"}
	} else {
		sortString = strings.Split(sort, ",")
	}

	logs.Beelog.Debug("query(%v)", util.S2Json(query))
	var tmps []map[string]interface{}
	if err := db.C(db.CN_BLOG).Find(query).Sort(sortString...).Skip((page - 1) * pageCount).Limit(pageCount).All(&tmps); err != nil {
		logs.Beelog.Error("query(%v)-->err(%v)", util.S2Json(query), err)
		return nil, err
	}
	var retmap map[string]interface{} = make(map[string]interface{})
	retmap["data"] = tmps
	if total {
		n, err := db.C(db.CN_BLOG).Find(query).Count()
		if err != nil {
			logs.Beelog.Error("query(%v)-->err(%v)", util.S2Json(query), err)
			return nil, err
		}
		retmap["total"] = n
	}

	return retmap, nil

}
