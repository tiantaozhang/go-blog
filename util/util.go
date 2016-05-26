package util

import (
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
