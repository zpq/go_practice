package main

import (
	"fmt"
)

func fbi(num int, ch chan int) {
	if num == 0 {
		ch <- 0
	} else if num == 1 || num == 2 {
		ch <- 1
	} else {
		cl := make(chan int)
		go fbi(num-1, cl)
		fmt.Println("cl")
		tmp := <-cl
		cr := make(chan int)
		go fbi(num-2, cr)
		fmt.Println("cr")
		tmp += <-cr
		ch <- tmp
	}
}
func main() {
	ch := make(chan int)
	go fbi(8, ch)
	fmt.Println(<-ch)
}
