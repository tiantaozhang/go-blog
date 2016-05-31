package util

import (
	"encoding/json"
	"time"
)

func TimeN() int64 {
	time.Now().Local().UnixNano()
}

func TimeM() int64 {
	time.Now().Local().UnixNano() / 1e6
}

func TimeS() {
	time.Now().Local().UnixNano() / 1e9
}

func S2Json(s interface{}) string {
	if s == nil {
		return ""
	}
	js, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(js)
}
