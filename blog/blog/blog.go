package blog

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tiantaozhang/go-blog/blog/db"
	"github.com/tiantaozhang/go-blog/logs"
	"gopkg.in/mgo.v2/bson"
)

func NewBlogG(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func NewBlogP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func EditBlogG(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func EditBlogP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func IndexBlog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func View(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bid := ps.ByName("id")
	logs.Beelog.Debug("view bid:(%v)", bid)
	t := template.ParseFiles("../../view/view.tpl")
	t.Execute(w, nil)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	bid := ps.ByName("id")
	if err := db.C(db.CN_BLOG).Update(bson.M{"_id": bid}, bson.M{"status": BS_DELETE}); err != nil {
		logs.Beelog.Error("deleteblog bid(%v)-->err(%v)", bid, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write([]byte("DEl OK"))
	return
}
