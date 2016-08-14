package blog

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tiantaozhang/go-blog/util"

	"github.com/julienschmidt/httprouter"
	"github.com/tiantaozhang/go-blog/blog/db"
	"github.com/tiantaozhang/go-blog/logs"
	"gopkg.in/mgo.v2/bson"
)

func Unescaped(x string) interface{} { return template.HTML(x) }

func NewBlogG(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Debug("---------new blog get----------")
	t, err := template.ParseFiles("view/new.tpl")
	if err != nil {
		logs.Beelog.Error("view,parsefile(%v),err(%v)", "view/new.tpl", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}

func NewBlogP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Debug("---------new blog post----------")
	var blog Blog
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&blog); err != nil {
		logs.Beelog.Error("decode r.body err:(%v)", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, bid, err := db.NewBid()
	if err != nil {
		logs.Beelog.Error("db newbid err(%v)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blog.Id = bid
	blog.Time = util.TimeS()
	blog.Last = util.TimeS()
	blog.Status = BS_NORMAL
	blog.Views = 0
	//blog.Author
	_, err = db.C(db.CN_BLOG).Upsert(bson.M{"_id": blog.Id}, bson.M{"$setOnInsert": blog})
	if err != nil {
		logs.Beelog.Error("upsert blog(%v)-->error", util.S2Json(blog), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":0,"data":{"bid":"%v"}}`, bid)))
	return

}

func EditBlogG(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Debug("---------edit blog get----------")
	id := r.PostFormValue("id")
	m, err := ListBlog("", "", 1, 1, id, "", nil, nil, "", false, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("view/edit.tpl")
	if err != nil {
		logs.Beelog.Error("edit,parsefile(%v),err(%v)", "view/new.tpl", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blog := map[string]interface{}{}
	if m["data"] != nil {

		if m["data"].([]map[string]interface{}) != nil {
			blog = m["data"].([]map[string]interface{})[0]
		}
	}
	t.Execute(w, blog)

}

func EditBlogP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Debug("---------edit blog post----------")
	if err := r.ParseForm(); err != nil {
		logs.Beelog.Error("parseform err(%v)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	update := bson.M{}
	title := r.Form.Get("title")
	id := r.Form.Get("id")
	t := r.Form.Get("type")
	if id == "" {
		logs.Beelog.Error("id nil err(%v)")
		http.Error(w, "id is nil", http.StatusForbidden)
		return
	}
	if title != "" {
		update["title"] = title
	}
	if t != "" {
		update["type"] = t
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Beelog.Error("read body err(%v)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if string(content) != "" {
		update["content"] = string(content)
	}

	if update == nil {
		logs.Beelog.Error("nothing to update")
		w.Write([]byte("nothing to update"))
		return
	}
	if err := db.C(db.CN_BLOG).Update(bson.M{"_id": id}, bson.M{"$set": update}); err != nil {
		logs.Beelog.Error("query(%v),update(%v)-->err(%v)", util.S2Json(bson.M{"_id": id}), util.S2Json(bson.M{"$set": update}), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("save success!"))
	return
}

func IndexBlog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Debug("---------Index blog get----------")

	t := template.New("index")
	t = t.Funcs(template.FuncMap{"unescaped": Unescaped})
	t, err := template.ParseFiles("view/index.tpl")

	if err != nil {
		logs.Beelog.Error("index blog get parsefiles(%v) err(%v)", "view/index.tpl", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//blogs := []map[string]map[string]interface{}{}
	blog := map[string]interface{}{}
	// type Index struct {
	// 	Blogs []map[string]interface{} `json:"blogs"`
	// 	Title string                   `json:"title"`
	// }
	// var index Index
	uid := ""
	m, err := ListBlog(uid, "", 1, 20, "", "", nil, nil, "-time,-views", false, true)
	if err != nil {
		logs.Beelog.Error("listblog err:(%v)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blogs := []map[string]interface{}{}

	blogs = m["data"].([]map[string]interface{})

	blog["blogs"] = blogs
	blog["title"] = "what the fuck"
	blog["total"] = m["total"]
	//index.Title = "what the fuck"
	logs.Beelog.Debug("index:(%v)", util.S2Json(blog))

	t.Execute(w, blog)
}

func View(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bid := ps.ByName("id")
	logs.Beelog.Debug("view bid:(%v)", bid)
	if bid == "" {
		logs.Beelog.Error("bid is nil")
		http.Error(w, "bid is nil", http.StatusForbidden)
		return
	}
	t, err := template.ParseFiles("view/view.tpl")
	if err != nil {
		logs.Beelog.Error("view,parsefile(%v) err(%v)", "../../view/view.tpl", err)
		if t, err = template.ParseFiles("../../view/view.tpl"); err != nil {
			logs.Beelog.Error("view,parsefile(%v),err(%v)", "../view/view.tpl", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var b Blog
	// bid = fmt.Sprintf("b%v", bid)
	if err = db.C(db.CN_BLOG).Find(bson.M{"_id": bid, "status": BS_NORMAL}).One(&b); err != nil {
		logs.Beelog.Error("query(%v)-->(%v)", util.S2Json(bson.M{"_id": bid, "status": BS_NORMAL}), err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tm := time.Unix(b.Time, 0)
	tmForString := tm.Format("2006-01-02 15:04:05 PM")

	m := map[string]interface{}{
		"title":   b.Title,
		"author":  b.Author,
		"content": b.Content,
		"time":    tmForString,
	}

	t.Execute(w, m)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Beelog.Warning("---------delete blog get----------")
	bid := ps.ByName("id")
	if err := db.C(db.CN_BLOG).Update(bson.M{"_id": bid}, bson.M{"$set": bson.M{"status": BS_DELETE}}); err != nil {
		logs.Beelog.Error("deleteblog bid(%v)-->err(%v)", bid, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		//http.Redirect(w, r, "/indexblog", http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, "/indexblog", http.StatusMovedPermanently)
	//	w.Write([]byte("DEl OK"))
	return
}
