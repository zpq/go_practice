package main

import (
	"log"
	"net"
	"runtime"
)

const (
	ip   = ""
	port = "8090"
)

func main() {
	// listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	runtime.GOMAXPROCS(4)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+port)
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("tcp listen falied!")
	}
	Server(listen)
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("client excrption : " + err.Error())
			continue
		}
		log.Println("a connection comes from remote client " + conn.RemoteAddr().String())
		defer conn.Close()
		go handleCLient(conn)
	}
}

func handleCLient(conn net.Conn) {
	data := make([]byte, 128)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println("a error happened when read data from client " + err.Error())
			break
		}
		log.Print(string(data[0:n]))
	}
}
