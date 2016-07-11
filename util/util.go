package util

import (
	"encoding/json"
	"time"
	"log"
)

func TimeN() int64 {
	return time.Now().Local().UnixNano()
}

func TimeM() int64 {
	return time.Now().Local().UnixNano() / 1e6
}

func TimeS()int64 {
	return  time.Now().Local().UnixNano() / 1e9
}

func S2Json(s interface{}) string {
	if s == nil {
		log.Fatal("s(%v) is nil",s )
		return ""
	}
	js, err := json.Marshal(s)
	if err != nil {
		log.Fatal("marshal s(%v) -->error(%v)",s,err )
		return ""
	}
	return string(js)
}
