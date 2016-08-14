package blog

import (
	"strings"
	"github.com/tiantaozhang/go-blog/logs"
	"github.com/tiantaozhang/go-blog/blog/db"
	"testing"
	"github.com/tiantaozhang/go-blog/util"
	"gopkg.in/mgo.v2/bson"

)

func TestList(t *testing.T) {

	if _,err:=ListBlog("", "", 0, 0, "", "title", []string{"tags"}, []int{BT_ORIGIN}, "", false, true);err==nil{
		t.Error("err1")
		return
	}
	if m,err:=ListBlog("", "", 1, 15, "", "", nil, []int{BT_ORIGIN}, "", false, true);err!=nil {
		t.Error(err)
		return
	}else{
		logs.Beelog.Debug(" listblog  total:(%v) data(%v)", m["total"],m["data"])
	}

	if m,err:=ListBlog("", "", 1, 15, "", "", nil, []int{BT_ORIGIN}, "-time,-views", false, true);err!=nil {
		t.Error(err)
		return
	}else{
		logs.Beelog.Debug(" listblog  total:(%v) data(%v)", m["total"],m["data"])
	}
	


}

func Remove()  {
	db.C(db.CN_BLOG).RemoveAll(nil)
}


func init() {
	Remove()
	add("u1","tatumn","the 1th blog", BT_ORIGIN, "第一篇,1111",0)
	add("u2","tatumn","the 2th blog", BT_ORIGIN, "第二篇,2222",1)
	add("u3","tatumn","the 3th blog", BT_ORIGIN, "第三篇,1111",2)
	add("u4","tatumn","the 4th blog", BT_ORIGIN, "第四篇,2222",3)


}


func add(uid,author ,title string,btype int,tags string,views int ) {
	var blog Blog
	_, bid, _ := db.NewBid()
	blog.Id = bid
	blog.Time = util.TimeS()
	blog.Last = util.TimeS()
	blog.Status = BS_NORMAL
	blog.Views = views

	blog.Title=title
	blog.Type=btype
	blog.Tags=strings.Split(tags, ",")
	blog.Author=author
	//blog.Author
	_, err := db.C(db.CN_BLOG).Upsert(bson.M{"_id": blog.Id}, bson.M{"$setOnInsert": blog})
	if err != nil {
		logs.Beelog.Error("upsert blog(%v)-->error", util.S2Json(blog), err)
		return
	}
}