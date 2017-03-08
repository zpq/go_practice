package main

import (
	"time"
)

type Item struct {
	v          interface{}
	expiration int64
}

func (item *Item) Expire() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}
