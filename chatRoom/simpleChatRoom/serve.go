package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	users           map[string]*User
	clients         map[string]*websocket.Conn
	refreshUserList chan bool
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type User struct {
	Name      string `json:"name"`
	LastAlive int64  `json:"lastAlive"`
}

type Message struct {
	Content string `json:"content"`
}

type Response struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

//delete
func GcOfflineUser() {

}

func refreshUserListsHandle() {
	for {
		refresh := <-refreshUserList
		if refresh {
			res := Response{2, "refresh Userlists", nil}
			for _, v := range users {
				res.Data = append(res.Data, v)
			}
			for _, v := range clients {
				v.WriteJSON(res)
			}
		}
	}
}

func echoHandle(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {

}

func Ws(w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)
	// defer conn.Close()
	// if err != nil {
	// 	log.Println(err.Error())
	// } else {
	// 	go echoHandle(w, r, conn)
	// }

	go func() {
		conn, err := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
		c, err := r.Cookie("username")
		if err != nil {
			log.Println(err.Error())
			return
		}
		username := c.Value
		for {
			_, content, err := conn.ReadMessage()
			if err != nil {
				log.Println(err.Error())
				return
			}

			_, ok := users[username]
			if !ok {
				log.Println("invalid user")
			}

			_, ok = clients[username]
			if !ok {
				clients[username] = conn
			}

			// log.Println(clients)
			var writeMessage string
			if string(content) == "hello!" {
				writeMessage = username + " enter the chat room! "
			} else {
				writeMessage = username + " says: " + string(content)
			}
			msg := Message{writeMessage}
			res := Response{1, "send message success", nil}
			res.Data = append(res.Data, msg)
			for _, v := range clients {
				err = v.WriteJSON(res)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()

}

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err.Error())
	}
	t.Execute(w, "ws://"+r.Host+"/Ws")
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	log.Println(username, "asd")
	res := Response{0, "username empty", nil}
	if username != "" {
		_, ok := users[username]
		if !ok {
			users[username] = &User{username, time.Now().Unix()}
			res.Data = append(res.Data, users[username])
			res.Status = 1
			res.Message = "login success"
			expiration := time.Now()
			expiration = expiration.AddDate(1, 0, 0)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			refreshUserList <- true
		} else {
			res.Message = "username exists"
		}
	}
	body, _ := json.Marshal(res)
	w.Write([]byte(body))
}

func main() {
	users = make(map[string]*User)
	clients = make(map[string]*websocket.Conn)
	refreshUserList = make(chan bool)
	go refreshUserListsHandle()
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/ws", Ws)
	http.ListenAndServe(":8008", nil)

}
