package main

import (
	"fmt"
)

const (
	maxTurn = 10
)

type Battle struct {
	ID             string
	BattleCategory int    // 战斗类型 pve or pvp
	RaidID         string // 副本id  适用于pve
	RaidCategory   int    // 副本类型 使用于pve
	ScenceID       string
	Groups         []*Group
	IsBattleGegin  bool
	IsBattleEnd    bool
	IsSetUpEnd     bool // 初始化是否已经结束
	Users          []*User
	MaxTurn        int // 战斗最大回合数
}

func (b *Battle) SetUpPVP(u1, u2 *User) {
	b.Users = []*User{u1, u2}
	b.Groups = append(b.Groups, &Group{UserID: u1.ID}, &Group{UserID: u2.ID})
	b.Groups[0].SetDecks(b.Users[0].DefaultDeck)
	b.Groups[0].SetHero(b.Users[0].DefaultHero)
	b.Groups[1] = &Group{}
	b.Groups[1].SetDecks(b.Users[1].DefaultDeck)
	b.Groups[1].SetHero(b.Users[1].DefaultHero)
	b.MaxTurn = maxTurn
	b.IsSetUpEnd = true
	b.IsBattleGegin = true

	fmt.Printf("Battle setup ends!\n")
	fmt.Printf("Battle id %s begin!\n", b.ID)
}

func (b *Battle) setUpPVE() {

}
