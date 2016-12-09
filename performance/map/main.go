package main

import (
	"log"
	"time"
)

func main() {
	m := make(map[int]int)
	l := 100000
	for i := 0; i < l; i++ {
		m[i] = i
	}
	log.Println("down")

	t1 := time.Now()
	log.Println("begin")
	nt := time.Now().UnixNano()
	for i := 0; i < l; i++ {
		_, ok := m[i]
		if ok {

		}
	}
	log.Println("end ", time.Now().Sub(t1), " -- ", time.Now().UnixNano()-nt)

	nt = time.Now().UnixNano()
	for _, v := range m {
		for _, vv := range m {
			if v == vv {
				break
			}
		}
	}
	log.Println("end ", time.Now().Sub(t1), " -- ", time.Now().UnixNano()-nt)
}
