package main

type BattleControl struct {
	*Battle
	Groups       map[*User]*Group
	Phase        string       // 战斗的当前阶段： 未开始， 开始， 进行中， 结束
	CurrentGroup *Group       // 当前的行动方
	CurrentTurn  int          // 当前的回合数
	CurrentActor *GameBattler // 当前的行动者
	IsActorEnd   bool         // 某一具体的行动者的行动是否结束
	IsGroupEnd   bool         // 某一方行动是否结束
	IsTurnEnd    bool         // 某一回合是否结束，即双方各自执行了一次行动
}

// MakeOrder ...生成战斗顺序
func (b *BattleControl) MakeOrder() {

}

func (b *BattleControl) ChangeCurrentGroup() {
	for k, _ := range b.Groups {
		if b.CurrentGroup != b.Groups[k] {
			b.CurrentGroup = b.Groups[k]
		}
	}
}

func (b *BattleControl) ChangeCurrentActor() {

}

func (b *BattleControl) TurnAdd() {
	b.CurrentTurn++
}

func (b *BattleControl) TurnBegin() {

}

func (b *BattleControl) TurnEnd() {

}

func (b *BattleControl) ProsessTurn() {
	b.GroupTurn()
}

func (b *BattleControl) GroupTurn(cg *Group) {

}

func (b *BattleControl) ActorTurn(at *GameBattler) {

}
