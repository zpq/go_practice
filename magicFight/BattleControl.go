package main

import "math/rand"
import "fmt"

type BattleControl struct {
	*Battle
	Phase        string       // 战斗的当前阶段： 未开始， 开始， 进行中， 结束
	CurrentGroup *Group       // 当前的行动方
	TargetGroup  *Group       // 当前的目标方
	CurrentTurn  int          // 当前的回合数
	CurrentActor *GameBattler // 当前的行动者
	IsActorEnd   bool         // 某一具体的行动者的行动是否结束
	IsGroupEnd   bool         // 某一方行动是否结束
	IsTurnEnd    bool         // 某一回合是否结束，即双方各自执行了一次行动
}

// MakeOrder ...生成战斗顺序
func (b *BattleControl) MakeOrder() {
	firstGroupIndex := rand.Intn(len(b.Groups))
	b.CurrentGroup = b.Groups[firstGroupIndex]
	if firstGroupIndex == 0 {
		b.TargetGroup = b.Groups[1]
	} else {
		b.TargetGroup = b.Groups[0]
	}
}

func (b *BattleControl) ChangeCurrentGroup() {
	tmpGroup := b.CurrentGroup
	b.CurrentGroup = b.TargetGroup
	b.TargetGroup = tmpGroup
}

// ChangeCurrentActor ...改变当前活动的战斗单元（英雄或者卡牌）
func (b *BattleControl) ChangeCurrentActor(g *GameBattler) {
	b.CurrentActor = g
}

func (b *BattleControl) TurnAdd() {
	b.CurrentTurn++
}

// func (b *BattleControl) TurnBegin() {

// }

// func (b *BattleControl) TurnEnd() {

// }

// ProsessTurn ...战斗流程处理
func (b *BattleControl) ProsessTurn() {
	b.TurnAdd()
	for {
		if b.CheckBattleEnd() {
			break
		}
		b.ProcessGroupTurn()

		fmt.Printf("Turn %d end!\n", b.CurrentTurn)
		b.ChangeCurrentGroup()
		b.IsTurnEnd = true
		b.TurnAdd()
		b.CurrentGroup.MinusTurnCoolDown()
	}
	fmt.Println("battle end!")
}

// ProcessGroupTurn ...群组处理
func (b *BattleControl) ProcessGroupTurn() {
	b.CurrentGroup.SummonCard()
	b.ChangeCurrentActor(b.CurrentGroup.Hero) // 每回合开始总是英雄第一个行动
	for _, v := range b.CurrentGroup.DeckInBoard {
		fmt.Println(v.Name)
		b.ProcessActorTurn(v)
		b.ChangeCurrentActor(v)
	}
	b.IsGroupEnd = true
}

// ProcessActorTurn ...具体的某一个行动者处理
func (b *BattleControl) ProcessActorTurn(g *GameBattler) {
	g.ApplyBuff(g)
	//伤害结果检查
	g.UseSkill(g)
	//伤害结果检查
	g.CommonAttack(g)
	//伤害结果检查
	b.IsActorEnd = true
}

func (b *BattleControl) DamageCheck() {
	b.CurrentGroup.CheckDeadCardsAndRemoveThem()
	b.TargetGroup.CheckDeadCardsAndRemoveThem()
}

func (b *BattleControl) CheckBattleEnd() bool {
	if b.CurrentTurn > b.MaxTurn {
		return true
	}
	return false
}
