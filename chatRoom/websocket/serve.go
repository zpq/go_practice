package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8008", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients map[string]*websocket.Conn

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	for {
		messageType, content, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err.Error())
		}
		remoteAddr := conn.RemoteAddr().String()

		_, ok := clients[remoteAddr]
		if !ok {
			clients[remoteAddr] = conn
		}
		log.Printf("Remote addr is : %s", remoteAddr)
		log.Printf("%s says : %s", remoteAddr, string(content))
		log.Printf("messageType %d", messageType)
		log.Println(clients)
		var writeMessage string
		if string(content) == "hello!" {
			writeMessage = remoteAddr + " enter the chat room! "
		} else {
			writeMessage = remoteAddr + " says: " + string(content)
		}
		for _, v := range clients {
			err = v.WriteMessage(messageType, []byte(writeMessage))
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err.Error())
	}
	t.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	clients = make(map[string]*websocket.Conn)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(*addr, nil)
}
