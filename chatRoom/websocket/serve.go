package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

var clients map[string]*websocket.Conn
var users map[string]*User

var tokenEncodeString string = "helloworld"
var addr = flag.String("addr", "localhost:8008", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     myCheckOrigin,
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoomId   string `json:"roomId"`
	Token    string `json:"token"`
}

type Room struct {
	Name    string
	Cap     int
	UserNum int
	IsOpen  bool
}

type Res struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type LoginRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// allow cross domain request
func myCheckOrigin(req *http.Request) bool {
	// if req.Header.Get("Origin") != "http://"+req.Host || req.Header.Get("Origin") != "http://localhost:8088" {
	// 	return false
	// }
	return true
}

func echo(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.Host)
	// log.Println(r.Header.Get("Origin"))
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

func getUserLists(w http.ResponseWriter, r *http.Request) {

}

func getRoomLists(w http.ResponseWriter, r *http.Request) {

}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Add("Access-Control-Allow-Methods", "POST,OPTIONS")
	// w.Header().Set("content-type", "application/json") //返回数据格式是json
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	res := Res{}
	res.Status = 0
	res.Message = "failed"
	if username == "admin" && password == "admin" {
		res.Status = 1
		res.Message = "login success"
		lr := LoginRes{username, getToken(username)}
		res.Data = append(res.Data, lr)
	}
	body, _ := json.Marshal(res)
	w.Write(body)
}

func checkToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization")

	tokenString := r.Header.Get("Authorization")
	fmt.Println("token from client is : ", tokenString)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return tokenEncodeString, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}

func getToken(username string) string {
	// w.Header().Set("Access-Control-Allow-Methods", "*")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":      "bar",
		"username": username,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(tokenEncodeString))
	// Sign and get the complete encoded token as a string using the secret

	if err != nil {
		log.Println(err.Error())
		return ""
	} else {
		return tokenString
	}

}

func main() {
	clients = make(map[string]*websocket.Conn)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/login", login)

	//restful api
	mux := pat.New()
	mux.Post("/", http.HandlerFunc(home))
	mux.Get("/room/:id/users", http.HandlerFunc(getUserLists))
	mux.Get("/rooms", http.HandlerFunc(getRoomLists))
	http.Handle("/", mux)

	http.ListenAndServe(*addr, nil)
}
