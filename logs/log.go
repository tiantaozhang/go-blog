package logs

import (
	"github.com/astaxie/beego/logs"
)

var Beelog *logs.BeeLogger

func init() {
	Beelog = logs.NewLogger(0)
	Beelog.SetLogger("console", "")
	Beelog.EnableFuncCallDepth(true)
	Beelog.SetLevel(logs.LevelDebug)
	Beelog.SetLogFuncCallDepth(2)
}
