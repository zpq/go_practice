package main

import (
	"fmt"
	"time"
)

func timeout(t time.Duration) func() {
	start := time.Now()
	return func() {
		if time.Now().Sub(start) > t {
			fmt.Println("timeout")
		}
	}
}

func closureDataTest() {
	x := 0
	go func() {
		x++
		fmt.Println(x)
	}()
	x++
}

func main() {
	defer timeout(time.Second)()
	time.Sleep(time.Second * 2)
	fmt.Println("hello world!")
	closureDataTest()
	time.Sleep(time.Second * 10)
}
