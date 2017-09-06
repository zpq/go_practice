package main

type GameBattler struct {
	*GameBattlerBase
	// ActionList    []*GameAction
	// CurrentAction int
	*GameAction
}

// PerformAction ...战斗者执行动作 动作分为: 1.处理buff的效果 2.执行技能 3.执行普通攻击
func (g *GameBattler) PerformAction() {
}

// 攻击需要提供攻击类型， 攻击数值
//（这样的话, 方守方可以根据攻击类型, 判断是否可以抵抗莫衷攻击类型, 计算时利用免伤buff, 来减少伤害）
