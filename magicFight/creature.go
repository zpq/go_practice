package main

// Creature ...base model
type Creature struct {
	ID               string
	Name             string
	Level            int
	ActorType        int // 1是英雄 2是卡牌
	Hp               int
	Attack           int
	Rarity           int // 1 common 2 elite 3 eqic 4 legend
	Faction          int // 1 human  2 Zerg  3 elf  4 Undead 5 protoss
	Star             int
	Equips           []*Equipment
	Skills           []*Skill
	Buffs            []*Buff
	initTurnCooldown int
}

// func (c *Creature) GetFields() *Creature {
// 	return c
// }
