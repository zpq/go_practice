package main

type GameAction struct {
	ID   string
	Name string
}

func (g *GameAction) Excute(gb *GameBattler) {
	g.ApplyBuff(gb)
	g.UseSkill(gb)
	g.CommonAttack(gb)
}

// CommonAttack ...普通攻击 脚本实现
func (g *GameAction) CommonAttack(gb *GameBattler) {

}

// UseSkill ...使用技能 脚本实现
func (g *GameAction) UseSkill(gb *GameBattler) {

}

// ApplyBuff ...buff效果实施 脚本实现
func (g *GameAction) ApplyBuff(gb *GameBattler) {

}
