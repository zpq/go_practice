package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	rooms            []*Room
	users            map[string]*User
	clients          map[string]*websocket.Conn
	serverhost       string = ":8008"
	validRemoteHosts string = "localhost:8088"
	upgrader                = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     myCheckOrigin,
	}
)

type Room struct {
	id      int
	members []*User
	guests  []*User // future todo
	isWait  bool    //是否在等待另一个用户进入
	battle  Battle
}

type Battle struct {
	turn        string      //下一次请求应该是谁的，不符合的认定为非法请求，不予处理(使用user的token对的区分)
	Weather     int         `json:"weather"`
	BattleScore BattleScore `json:"battleScore"`
}

type BattleScore struct {
}

type User struct {
	Id           int          `json:"id"`
	Name         string       `json:"name"`
	CardInfo     CardInfo     `json:"cardInfo"`
	FightHistory FightHistory `json:"fightHistory"`
}

type CardInfo struct {
	TotalCards    []*Card `json:"totalCards"`
	UsedCards     []*Card `json:"usedCards"`
	UnUsedCards   []*Card `json:"unUsedCards"`
	InfantryCards []*Card `json:"infantryCards"` // active card
	ArcherCards   []*Card `json:"archerCards"`   // active card
	SlingCards    []*Card `json:"slingCards"`    // active card
	TotalDamage   int     `json:"totalDamage"`
}

type FightHistory struct { // 2-0 => 2;  2-1=>1; 1-2 => 0; 0-2 => -1  (0-2 common happened in run away)
	Score       int `json:"score"`
	Win         int `json:"win"`
	Lost        int `json:"lost"`
	Last        int `json:"last"` // record last PK result
	ContinueWin int `json:"continueWin"`
}

type Card struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	IsHero        bool   `json:"jsHero"`        //金卡
	IsSpy         bool   `json:"isSpy"`         //间谍卡
	GroupType     int    `json:"groupType"`     //卡组类型
	UnitType      int    `json:"unitType"`      //卡牌类型 0:weather 1:infantry 2:archer 3:sling
	WeatherEffect int    `json:"weatherEffect"` //天气类型 0:sun 1:debuff infantry 2:debuff archer  3:debuff sling
	BufferEffect  int    `json:"bufferEffect"`  //自带buff (maybe use callfunc)
	BaseDamage    int    `json:"baseDamage"`
	ComputeDamage int    `json:"computeDamage"`
	IsUsed        bool   `json:"isUsed"`
	IsActive      bool   `json:"isActive"` //是否正在被使用
}

func myCheckOrigin(r *http.Request) bool {
	if r.Host == validRemoteHosts {
		return true
	}
	return false
}

func Ws(w http.ResponseWriter, r *http.Request) {

}

func main() {

}
