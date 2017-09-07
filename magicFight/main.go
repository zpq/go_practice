package main

import (
	"log"
	"strconv"
	"strings"
)

var maxBattleID = 0

const (
	scriptPrefix      = "script/"
	scriptSkillPrefix = scriptPrefix + "skill/"
	scriptBuffPrefix  = scriptPrefix + "buff/"
)

var (
	cardSource  = make(map[string]Card)
	skillSource = make(map[string]Skill)
)

func main() {

	LoadSkillSource("./data/skill.txt")
	LoadCardSource("./data/card.txt")

	u1 := &User{ID: 1, Name: "peter"}
	u2 := &User{ID: 2, Name: "mary"}
	u1.MakeOneDeck()
	u1.MakeOneHero()
	u2.MakeOneDeck()
	u2.MakeOneHero()

	b := &Battle{ID: "battle1"}
	b.SetUpPVP(u1, u2)

	bc := &BattleControl{Battle: b}
	bc.MakeOrder()
	bc.ProsessTurn()

}

func LoadCardSource(filename string) {
	r, err := LoadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	for k, v := range r {
		if k == 0 {
			continue
		} else {
			c := Card{Creature: &Creature{}}
			c.ID = v[0]
			c.Name = v[1]
			c.Level = stringToInt(v[2])
			c.ActorType = stringToInt(v[3])
			c.Hp = stringToInt(v[4])
			c.Attack = stringToInt(v[5])
			c.Rarity = stringToInt(v[6])
			c.Faction = stringToInt(v[7])
			c.Star = stringToInt(v[8])
			c.initTurnCooldown = stringToInt(v[9])
			skills := strings.Split(v[10], ",")
			for _, sv := range skills {
				// if sv != "0" {
				skill := skillSource["skill"+sv]
				c.Skills = append(c.Skills, &skill)
				// }
			}
			cardSource[c.ID] = c
		}
	}
}

func LoadSkillSource(filename string) {
	r, err := LoadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	for k, v := range r {
		if k == 0 {
			continue
		} else {
			s := Skill{}
			s.ID = v[0]
			s.Name = v[1]
			s.ScriptID = v[2]
			s.Description = v[3]
			skillSource[s.ID] = s
		}
	}
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return i
}
