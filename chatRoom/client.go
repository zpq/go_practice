package main

import (
	"log"
	"net"
	"runtime"
	"time"
)

const (
	address = "127.0.0.1:8090"
)

func main() {
	runtime.GOMAXPROCS(4)
	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp4", address)
		if err != nil {
			log.Fatal("dial tcp failed : " + err.Error())
		}
		defer conn.Close()
		go Client(conn)
	}
}

func Client(conn net.Conn) {
	for i := 0; i < 1; i++ {
		data := make([]byte, 128)
		data = []byte("HelloWorld")
		conn.Write(data)
		time.Sleep(time.Millisecond * 100)
	}
}
