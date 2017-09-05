package main

const (
	maxTurn = 20
)

type Battle struct {
	ID             string
	BattleCategory int // 战斗类型 pve or pvp
	RaidID         int // 副本id
	RaidCategory   int // 副本类型
	IsBattleGegin  bool
	IsBattleEnd    bool
	Users          []*User
	MaxTurn        int // 战斗最大回合数
}

func (b *Battle) SetUp() {

}
