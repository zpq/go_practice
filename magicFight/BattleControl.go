package main

import "math/rand"
import "fmt"
import "strconv"

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
	Logger       *BattleLog   // 战斗日志记录
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
		b.ShowStatus()
		b.ProcessGroupTurn()
		b.CurrentGroup.MinusTurnCoolDown()
		b.ChangeCurrentGroup()
		b.ProcessGroupTurn()
		b.CurrentGroup.MinusTurnCoolDown()
		b.ChangeCurrentGroup()
		fmt.Printf("Turn %d end!\n", b.CurrentTurn)
		b.IsTurnEnd = true
		b.TurnAdd()
	}
	fmt.Println("battle end!")
}

// ProcessGroupTurn ...群组处理
func (b *BattleControl) ProcessGroupTurn() {
	b.CurrentGroup.SummonCard()
	b.ShowStatus()
	b.ChangeCurrentActor(b.CurrentGroup.Hero)

	// hero actor
	b.ProcessActorTurn(b.CurrentGroup.Hero)
	b.ShowStatus()

	// card actor
	for _, v := range b.CurrentGroup.DeckInBoard {
		b.ProcessActorTurn(v)
		b.ChangeCurrentActor(v)
	}
	b.IsGroupEnd = true
}

// ProcessActorTurn ...具体的某一个行动者处理
func (b *BattleControl) ProcessActorTurn(g *GameBattler) {
	g.ApplyBuff(g, b)
	// b.DamageCheck()
	g.UseSkill(g, b)
	b.DamageCheck()
	g.CommonAttack(g, b)
	b.DamageCheck()
	b.IsActorEnd = true
}

func (b *BattleControl) DamageCheck() {
	b.CurrentGroup.CheckDeadCardsAndRemoveThem()
	b.TargetGroup.CheckDeadCardsAndRemoveThem()
	b.ShowStatus()
}

func (b *BattleControl) ShowStatus() {
	fmt.Println("\n--------------------show status start-------------------\n")
	fmt.Print("User" + strconv.Itoa(b.CurrentGroup.UserID) + " ")
	fmt.Print(b.CurrentGroup.Hero.Name + "(")
	fmt.Print(b.CurrentGroup.Hero.CAttack)
	fmt.Print(",")
	fmt.Print(b.CurrentGroup.Hero.CHP)
	fmt.Print(") ")
	for _, v := range b.CurrentGroup.DeckInBoard {
		fmt.Print(v.Name + "(")
		fmt.Print(v.CAttack)
		fmt.Print(",")
		fmt.Print(v.CHP)
		fmt.Print(") ")
	}
	// ------------------------------------------------ //
	fmt.Print("\nUser" + strconv.Itoa(b.TargetGroup.UserID) + " ")
	fmt.Print(b.TargetGroup.Hero.Name + "(")
	fmt.Print(b.TargetGroup.Hero.CAttack)
	fmt.Print(",")
	fmt.Print(b.TargetGroup.Hero.CHP)
	fmt.Print(") ")
	for _, v := range b.TargetGroup.DeckInBoard {
		fmt.Print(v.Name + "(")
		fmt.Print(v.CAttack)
		fmt.Print(",")
		fmt.Print(v.CHP)
		fmt.Print(") ")
	}
	fmt.Println("\n\n--------------------show status end-------------------\n")
}

func (b *BattleControl) CheckBattleEnd() bool {
	if b.CurrentTurn > b.MaxTurn {
		return true
	}
	return false
}
