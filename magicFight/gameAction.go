package main

import (
	"fmt"
	"log"
	"strconv"
)

type GameAction struct {
	ID   string
	Name string
}

func (g *GameAction) Excute(gb *GameBattler) {
	// g.ApplyBuff(gb)
	// g.UseSkill(gb)
	// g.CommonAttack(gb)
}

// CommonAttack ...普通攻击 脚本实现
func (g *GameAction) CommonAttack(gb *GameBattler, b *BattleControl) {
	for _, v := range gb.Skills {
		if v.ID != "skill0" {
			continue
		}
		fmt.Print("User")
		fmt.Print(b.CurrentGroup.UserID)
		fmt.Println(" " + gb.Name + " position" + strconv.Itoa(gb.Position) + " use skill " + v.Name)
		err, result := v.Use(gb, b)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(result)
	}
}

// UseSkill ...使用技能 脚本实现
func (g *GameAction) UseSkill(gb *GameBattler, b *BattleControl) {
	for _, v := range gb.Skills {
		if v.ID == "skill0" {
			continue
		}
		fmt.Print("User")
		fmt.Print(b.CurrentGroup.UserID)
		fmt.Println(" " + gb.Name + " position" + strconv.Itoa(gb.Position) + " use skill " + v.Name)
		err, result := v.Use(gb, b)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(result)
	}

}

// ApplyBuff ...buff效果实施 脚本实现
func (g *GameAction) ApplyBuff(gb *GameBattler, b *BattleControl) {

}
