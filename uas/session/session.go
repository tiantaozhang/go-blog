package session

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

type ManagerSession struct {
	Lock sync.RWMutex
	SM   map[interface{}]*list.Element
	SL   *list.List
	Name string
}

type Session struct {
	Lock    sync.RWMutex
	Sid     string
	Key     string
	Values  map[interface{}]interface{}
	Expires int
}
