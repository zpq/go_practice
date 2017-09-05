package main

type GameBattler struct {
	*GameBattlerBase
	Member GameEntity // new GameBattlerBase
	// ActionList    []*GameAction
	// CurrentAction *GameAction
}

func (g *GameBattler) PerformAction() {

}
