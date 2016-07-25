package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	StartClient()
}

func ClientSend(conn net.Conn) {
	var input string
	clientName := conn.LocalAddr().String()
	for {
		fmt.Scanln(&input)
		if input == "/q" || input == "/Q" {
			log.Fatal("user exit program forcely!")
		}
		_, err := conn.Write([]byte(clientName + " Says:: " + input))
		if err != nil {
			fmt.Println("client write error : " + err.Error())
			break
		}
	}
}

func ClientReciver(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			conn.Close()
			log.Fatal("server make boo boo!")
		}
		fmt.Println(string(buff[0:n]))
	}
}

func StartClient() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:3006")
	if err != nil {
		log.Fatal("make tcp ip failed " + err.Error())
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal("connect server failded " + err.Error())
	}

	go ClientSend(conn)

	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			conn.Close()
			log.Fatal("server make boo boo!")
		}
		fmt.Println(string(buff[0:n]))
	}
}
