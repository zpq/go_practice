package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	allCards map[int][]*Card
	rooms    map[int]*Room
	users    map[string]*User
	clients  map[string]*websocket.Conn
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     myCheckOrigin,
	}
	writeSignal chan string
	roomSignal  chan int
)

const (
	tokenSecret          string = "kuntCards_kuntCards_secrets"
	serverhost           string = ":8008"
	validRemoteHosts     string = "localhost:8088"
	reconnectionDeadline        = 60
	errorResType                = 0 // parse websocket message type
	successResType              = 1
	messageResType              = 2
	closeResType                = 8
)

type Response struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Datas   []interface{} `json:"datas"`
}

type Room struct {
	id          int
	members     []string
	guests      []string // future todo
	membersChan map[int]chan string
	isWait      bool // 是否在等待另一个用户进入
	isBegin     bool
	Battle      Battle `json:"battle"`
}

type Battle struct {
	cuser       string //当前的请求应该是谁的，不符合的认定为非法请求，不予处理(使用user的token来区分)
	turn        int
	Weather     []int          `json:"weather"`
	BattleScore map[string]int `json:"battleScore"` // exp:[username] => 2
}

type User struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	password     string
	token        string
	RoomId       int          `json:"roomId"`
	LastAlive    int64        `json:"lastAlive"`
	IsOnline     bool         `json:"isOnLine"`
	CardInfo     CardInfo     `json:"cardInfo"`
	FightHistory FightHistory `json:"fightHistory"`
}

type CardInfo struct {
	TotalCards  []*Card `json:"totalCards"`
	UsedCards   []*Card `json:"usedCards"`
	UnUsedCards []*Card `json:"unUsedCards"` // card which can use
	ActiveCards []*Card `json:"activeCards"`
	// InfantryCards []*Card       `json:"infantryCards"` // active card
	// ArcherCards   []*Card       `json:"archerCards"`   // active card
	// SlingCards    []*Card       `json:"slingCards"`    // active card
	WeatherCards []*Card       `json:"weatherCards"`
	BufferCards  map[int]*Card `json:"bufferCards"` // key => card position
	TotalDamage  int           `json:"totalDamage"`
}

type FightHistory struct { // 2-0 => 2;  2-1=>1; 1-2 => 0; 0-2 => -1  (0-2 common happened in run away)
	Score       int `json:"score"`
	Win         int `json:"win"`
	Lost        int `json:"lost"`
	Last        int `json:"last"` // record last PK result
	ContinueWin int `json:"continueWin"`
}

type Card struct {
	Id            int          `json:"id"`  // unique identify
	Cid           int          `json:"cid"` //
	Name          string       `json:"name"`
	IsHero        bool         `json:"jsHero"`        //英雄卡
	IsSpy         bool         `json:"isSpy"`         //间谍卡
	GroupType     int          `json:"groupType"`     //卡组类型
	UnitType      int          `json:"unitType"`      //卡牌类型 0:weather 1:infantry 2:archer 3:sling 4:buffer
	BrotherId     int          `json:"brotherId"`     //兄弟卡，一起出现有奇效！
	WeatherEffect int          `json:"weatherEffect"` //天气类型 0:sun 1:debuff infantry 2:debuff archer  3:debuff sling
	bufferEffect  BufferEffect //自带buff (use callback)
	BaseDamage    int          `json:"baseDamage"`
	ComputeDamage int          `json:"computeDamage"`
	IsUsed        bool         `json:"isUsed"`
	IsActive      bool         `json:"isActive"` //是否正在被使用
	Description   string       `json:"description"`
}

type BufferEffect func(card *Card, room *Room) error

func myCheckOrigin(r *http.Request) bool {
	if r.Host == validRemoteHosts {
		return true
	}
	return false
}

func Ws(w http.ResponseWriter, r *http.Request) {
	var isNew bool
	var token string
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error : " + err.Error())
		return
	}

	defer func() {
		log.Println("leave")
		// delete(clients, token) // important
		res := Response{
			Status:  8,
			Message: "server close connecting",
			Datas:   nil,
		}
		body, err := json.Marshal(res)
		if err == nil {
			conn.WriteControl(websocket.CloseMessage, body, time.Time{})
		}
		conn.Close()
	}()

	c, err := r.Cookie("kunt-token")
	if err != nil { //no token, need login
		log.Println(err.Error())
		return
	}

	token = c.Value
	tc, err := checkToken(token)
	if err != nil { // invalid token, need login
		log.Println(err.Error())
		return
	}

	username := tc["username"].(string)
	user := users[token]
	if user.Name != username {
		log.Println("username not match between token and server stored")
		return
	}

	// check whether user has a peice of cardGroup

	_, ok := clients[token]
	if !ok {
		clients[token] = conn
		users[token].LastAlive = time.Now().Unix()
		isNew = true //need dispatch room
	} else {
		isNew = false                                                // reconnect to room
		if time.Now().Unix()-user.LastAlive > reconnectionDeadline { // over maxTime
			delete(clients, token)
			return
		} else {
			clients[token] = conn
		}
	}

	go write(conn, isNew)

	if isNew {
		dispatchRoom(user, conn)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			break
		}
		user.LastAlive = time.Now().Unix()
		dataParse(string(msg)) // data parse

		writeSignal <- string(msg)
	}
}

func write(conn *websocket.Conn, isNew bool) {
	conn.WriteJSON(<-writeSignal)
}

func dataParse(msg string) {

}

func dispatchRoom(user *User, conn *websocket.Conn) {
	if len(rooms) <= 0 {
		createRoom(user)
	} else {
		isDone := false
		for _, v := range rooms {
			if v.isWait && len(v.members) == 1 {
				v.members = append(v.members, user.token)
				v.isWait = false
				v.isBegin = true
				// write some data
				isDone = true
				break
			}
		}
		if !isDone {
			createRoom(user)
		}
	}
}

func createRoom(user *User) *Room {
	room := &Room{}
	room.id = len(rooms) + 1
	room.isWait = true
	room.members = append(room.members, user.token)
	// write some data
	return room
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", validRemoteHosts)
	r.ParseForm()
	username := strings.Trim(r.PostFormValue("username"), "")
	password := strings.Trim(r.PostFormValue("password"), "")
	res := Response{0, "login failed, wrong username or password", nil}
	user, ok := checkUserExist(username)
	if ok {
		if password == user.password {
			token, err := createToken(username)
			if err != nil {
				log.Println(err.Error())
			} else {
				user.token = token
				user.LastAlive = time.Now().Unix()
				user.IsOnline = true
				res.Status = 1
				res.Message = "login success"
				res.Datas = append(res.Datas, user)
			}
		}
	}
	body, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write([]byte(body))
}

func Register(w http.ResponseWriter, r *http.Request) { // auto login after register
	w.Header().Set("Access-Control-Allow-Origin", validRemoteHosts)
	r.ParseForm()
	username := strings.Trim(r.PostFormValue("username"), "")
	password := strings.Trim(r.PostFormValue("password"), "")
	_, ok := checkUserExist(username)
	res := Response{0, "username already exists", nil}
	if !ok {
		user, token, err := createUser(username, password)
		if err != nil {
			user.LastAlive = time.Now().Unix()
			user.IsOnline = true
			res.Datas = append(res.Datas, user)
			expiration := time.Now()
			expiration = expiration.AddDate(365, 0, 0)
			cookie := &http.Cookie{
				Name:     "kunt-token",
				Value:    token,
				Expires:  expiration,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
	}
	body, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write([]byte(body))
}

func createUser(username, password string) (*User, string, error) {
	token, err := createToken(username)
	user := &User{Name: username, password: password}
	if err != nil {
		log.Println(err.Error())
	} else {
		user.Id = len(users) + 1
		user.token = token
		users[token] = user
	}
	return user, token, err
}

func checkUserExist(username string) (*User, bool) {
	for _, v := range users {
		if v.Name == username {
			return v, true
		}
	}
	return &User{}, false
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"Expires":  time.Now().Add(time.Second * 3600 * 24 * 30).Unix(),
		// "nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func checkToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

//user select cards
func selectCards(w http.ResponseWriter, r *http.Request) {
	// allCards map[int][]*Card
	w.Header().Set("Access-Control-Allow-Origin", validRemoteHosts)
	c, err := r.Cookie("kunt-token")
	if err != nil { // need login
		log.Println("no token exist")
	}

	token := c.Value
	tc, err := checkToken(token)
	if err != nil { //need login
		log.Println("invalid token")
	}

	username := tc["username"].(string)
	user := users[token]
	if user.Name != username { // need login
		log.Println("username not match between token and server stored")
	}

	r.ParseForm()
	cards := r.PostFormValue("cards")
	log.Printf("%T\n", cards)

	l := len(cards)
	for i := 0; i < l; i++ {

	}

}

//calculate points
func calPoints(room *Room) {

}

func deleteCardFromCardsById(cards []*Card, id int) {
	l := len(cards)
	var index int
	for i := 0; i < l; i++ {
		if id == cards[i].Id {
			index = i
		}
	}
	if l > 0 {
		cards = append(cards[:index], cards[index+1:]...)
	}
}

func getOpponent(room *Room) *User {
	for i := 0; i < len(room.members); i++ {
		if room.members[i] != room.Battle.cuser {
			return users[room.members[i]]
		}
	}
	return nil
}

func getRandomCardFromTotalCards(user *User) *Card {
	l := len(user.CardInfo.TotalCards)
	return user.CardInfo.TotalCards[rand.Intn(l)]
}

func clearWeatherCards(user *User) {
	l := len(user.CardInfo.WeatherCards)
	for i := 0; i < l; i++ {
		user.CardInfo.UsedCards = append(user.CardInfo.UsedCards, user.CardInfo.WeatherCards[i])
	}
	user.CardInfo.WeatherCards = user.CardInfo.WeatherCards[l:]
}

//api
func useCard(room *Room, card *Card) {
	token := room.Battle.cuser // sender
	user := users[token]
	if card.UnitType == 0 { // is weather card
		if card.WeatherEffect == 0 { // sum weather,clear weather card
			room.Battle.Weather = room.Battle.Weather[len(room.Battle.Weather):] // clear slice
			userOp := getOpponent(room)
			if userOp == nil {
				log.Fatal("Fatal error!No opponent exist!")
			}
			clearWeatherCards(user)
			clearWeatherCards(userOp)
		} else {
			room.Battle.Weather = append(room.Battle.Weather, card.Id)
			user.CardInfo.ActiveCards = append(user.CardInfo.ActiveCards, card)
			user.CardInfo.WeatherCards = append(user.CardInfo.WeatherCards, card)
			deleteCardFromCardsById(user.CardInfo.UnUsedCards, card.Id)
		}
	} else if card.UnitType == 4 { // is pury buff/debuff card
		user.CardInfo.UsedCards = append(user.CardInfo.UsedCards, card)
		deleteCardFromCardsById(user.CardInfo.UnUsedCards, card.Id)
		card.bufferEffect(card, room) // excute callback
	} else { // UnitType 1,2,3
		if card.IsSpy {
			deleteCardFromCardsById(user.CardInfo.UnUsedCards, card.Id)
			deleteCardFromCardsById(user.CardInfo.TotalCards, card.Id)
			c1, c2 := getRandomCardFromTotalCards(user), getRandomCardFromTotalCards(user)
			user.CardInfo.UnUsedCards = append(user.CardInfo.UnUsedCards, c1)
			user.CardInfo.UnUsedCards = append(user.CardInfo.UnUsedCards, c2)
			userOp := getOpponent(room)
			if userOp == nil {
				log.Fatal("Fatal error!No opponent exist!")
			}
			userOp.CardInfo.TotalCards = append(userOp.CardInfo.TotalCards, card)
			userOp.CardInfo.UnUsedCards = append(userOp.CardInfo.UnUsedCards, card)
		} else {
			deleteCardFromCardsById(user.CardInfo.UnUsedCards, card.Id)
			user.CardInfo.ActiveCards = append(user.CardInfo.ActiveCards, card)
			card.bufferEffect(card, room)
		}
	}
	calPoints(room)
}

func getAllCards(w http.ResponseWriter, r *http.Request) {

}

func battleBegin(room *Room) {
	//send battle start signal
	conn1, conn2 := clients[room.members[0]], clients[room.members[1]]
	go write(conn1, isNew)
	go write(conn2, isNew)
	for {

	}
}

func roomHandle() {
	for {
		roomId := <-roomSignal
		go battleBegin(rooms[roomId])
	}
}

func main() {

	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/getAllCards", getAllCards)
	http.HandleFunc("/selectCards", selectCards)
	if err := http.ListenAndServe(serverhost, nil); err != nil {
		log.Fatal(err.Error())
	}
}
