package logs

import (
	"github.com/astaxie/beego/logs"
)

var Beelog *logs.BeeLogger

func init() {
	beelog = logs.NewLogger(0)
	beelog.SetLogger("console", "")
	beelog.EnableFuncCallDepth(true)
	beelog.SetLevel(logs.LevelDebug)
	beelog.SetLogFuncCallDepth(2)
}
