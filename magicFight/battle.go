package main

const (
	maxTurn = 20
)

type Battle struct {
	ID             string
	BattleCategory int // 战斗类型 pve or pvp
	RaidID         int // 副本id
	RaidCategory   int // 副本类型
	IsGegin        bool
	IsEnd          bool
	Users          []*User
	Groups         map[*User]*Group
	CurrentGroup   *Group
	CurrentTurn    int    // 当前的回合数
	MaxTurn        int    // 战斗最大回合数
	IsOneTurnEnd   int    // 某一回合是否结束，即双方各自执行了一次行动
	Phase          string // 战斗的当前阶段： 未开始， 开始， 进行中， 结束
}

func (b *Battle) SetUp() {

}
