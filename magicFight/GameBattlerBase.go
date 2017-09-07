package main

type GameBattlerBase struct {
	*Creature
	CHP           int
	CAttack       int
	IsInBoard     bool
	IsInGraveYard bool
	Position      int
}

func NewGameBattlerBase() *GameBattlerBase {
	return &GameBattlerBase{
		Creature: &Creature{},
	}
}

func (g *GameBattler) CopyCard(c *Card) {
	g.ID = c.ID
	g.Name = c.Name
	g.ActorType = c.ActorType
	g.Faction = c.Faction
	g.Star = c.Star
	g.Hp = c.Hp
	g.Attack = c.Attack
	g.Rarity = c.Rarity
	g.initTurnCooldown = c.initTurnCooldown
	g.CHP = g.Hp
	g.CAttack = g.Attack
	// g.Buffs = c.Buffs
	g.Skills = c.Skills
	g.Equips = c.Equips
}

// IsDead ...check is dead
func (g *GameBattlerBase) IsDead() bool {
	return g.CHP <= 0
}

// // SetUpSkills ... 根据等级从数据中心加载技能
// func (g *GameBattlerBase) SetUpSkills() {

// }

func (g *GameBattlerBase) AddBuff(b *Buff) bool {
	for k, v := range g.Buffs {
		if v == b {
			if b.IsCanPly {
				g.Buffs[k].ply++
				return true
			}
			return false
		}
	}
	g.Buffs = append(g.Buffs, b)
	return true
}

func (g *GameBattlerBase) RemoveBuff(b *Buff) {

}

func (g *GameBattlerBase) ContainsBuff(b *Buff) {

}

func (g *GameBattlerBase) CheckBuffIsCollision(a *Buff, b *Buff) bool {
	return false
}

func (g *GameBattlerBase) CanMove() bool {
	// g.ContainsBuff()  // 检查角色是否包含限制行动的debuff(从数据中心那数据)
	return false
}

func (g *GameBattlerBase) MinusTurnCoolDown() {
	g.initTurnCooldown--
}
