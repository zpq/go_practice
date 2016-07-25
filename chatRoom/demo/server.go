package main

import (
	"fmt"
	"log"
	"net"
)

const (
	ip   = ""
	port = ":3006"
)

func main() {
	StartServer()
}

func ErrHandle(err error, msg string) {
	if err != nil {
		log.Fatal(msg + " " + err.Error())
	}
}

func StartServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+port)
	ErrHandle(err, "tcp ip ")
	listen, err := net.ListenTCP("tcp", tcpAddr)
	ErrHandle(err, "listen tcp ")

	conns := make(map[string]net.Conn)
	message := make(chan string)

	go ServerSendHandle(&conns, message)

	for {
		fmt.Println("server is listening...")
		conn, err := listen.Accept()
		ErrHandle(err, "accepting ")
		_, ok := conns[conn.RemoteAddr().String()]
		if !ok {
			conns[conn.RemoteAddr().String()] = conn
		}
		go ServerReceiveHandle(conn, message)
	}

}

func ServerReceiveHandle(conn net.Conn, message chan string) {
	buf := make([]byte, 1024)
	for {
		dl, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		message <- string(buf[0:dl])
	}
}

func ServerSendHandle(conns *map[string]net.Conn, message chan string) {
	for {
		msg := <-message
		fmt.Println("sending message to all clients: " + msg)
		for k, v := range *conns {
			_, err := v.Write([]byte(msg))
			if err != nil { //client is broken
				delete(*conns, k)
			}
		}
	}
}
