package blog

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/tiantaozhang/go-blog/blog/db"

	"docker/vendor/src/github.com/docker/go/canonical/json"
	"github.com/julienschmidt/httprouter"
	"github.com/tiantaozhang/go-blog/logs"
)

func TestRequest(t *testing.T) {

	route := httprouter.New()
	route.POST("/newblog", NewBlogP)
	route.POST("/editblog", EditBlogP)
	route.GET("/delete/:id", DeleteBlog)
	r1, _ := http.NewRequest("POST", "/newblog", strings.NewReader(`{"title":"1th blog","content":"this is the 1th blog,congratulation","type":10}`))
	w1 := httptest.NewRecorder()

	route.ServeHTTP(w1, r1)
	logs.Beelog.Debug("body:%v", w1.Body.String())
	var m map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &m)
	if err != nil {
		t.Error(err)
		return
	}
	data := m["data"]
	logs.Beelog.Debug("data:%v", data)
	m = m["data"].(map[string]interface{})
	bid := m["bid"]
	logs.Beelog.Debug("bid:%v", bid)

	r2, _ := http.NewRequest("POST", fmt.Sprintf("/editblog?id=%v&title=%v&type=", bid, "I am what Iam"), strings.NewReader("alter the 1th blog,congratulation"))
	w2 := httptest.NewRecorder()

	route.ServeHTTP(w2, r2)
	logs.Beelog.Debug("editblog:%v", w2.Body.String())

	r3, _ := http.NewRequest("POST", fmt.Sprintf("/editblog?id=%v&title=%v&type=%v", bid, "I am what Iam", 20), strings.NewReader("alter the 1th blog,congratulation"))
	w3 := httptest.NewRecorder()

	route.ServeHTTP(w3, r3)
	logs.Beelog.Debug("editblog:%v", w3.Body.String())

	r4, _ := http.NewRequest("GET", fmt.Sprintf("/delete/%v", bid), nil)
	w4 := httptest.NewRecorder()
	route.ServeHTTP(w4, r4)
	logs.Beelog.Debug("delete:%v", w4.Body.String())

}

func TestBinary(t *testing.T) {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		t.Error(err)
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Print(pi)
}
