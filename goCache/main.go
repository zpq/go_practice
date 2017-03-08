package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("begin.....")
	c := NewCache(time.Second * 3)
	expire, _ := time.ParseDuration("10s")
	c.Add("a", 123, expire)
	a, ok := c.Get("a")
	if !ok {
		fmt.Println("not found1")
		return
	}
	fmt.Println(a)
	time.Sleep(time.Second * 2)
	c.SaveToFile("./cache")
	time.Sleep(time.Second * 3)
	a, ok = c.Get("a")
	if !ok {
		fmt.Println("not found2")
	}
	fmt.Println(a)
	c.LoadFromFile("./cache")

	a, ok = c.Get("a")
	if !ok {
		fmt.Println("not found3")
		return
	}
	fmt.Println(a)

}
