package main

type GameBattlerBase struct {
	*Creature
	CHP     int
	CAttack int
}

// IsDead ...check is dead
func (g *GameBattlerBase) IsDead() bool {
	return g.GetFields().Hp == 0
}

// SetUpSkills ... 根据等级从数据中心加载技能
func (g *GameBattlerBase) SetUpSkills() {

}

func (g *GameBattlerBase) AddBuff(b *Buff) bool {
	for k, v := range g.Buff {
		if v == b {
			if b.IsCanPly {
				g.Buff[k].ply++
				return true
			}
			return false
		}
	}
	g.Buff = append(g.Buff, b)
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
