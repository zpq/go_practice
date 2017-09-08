package main

import (
	"fmt"
)

// Group ...群组对象
type Group struct {
	UserID          int
	Hero            *GameBattler
	Deck            []*GameBattler
	DeckInBoard     []*GameBattler
	DeckInGraveYard []*GameBattler
}

// IsDeckInHandEmpty ...检查手牌是否为空
func (g *Group) IsDeckInHandEmpty() bool {
	return len(g.Deck) == 0
}

// IsDeckInBoardEmpty ...检查场上是否存在卡牌
func (g *Group) IsDeckInBoardEmpty() bool {
	return len(g.DeckInBoard) == 0
}

// Clear ...清除群组中的所有卡牌
func (g *Group) Clear() {
	g.Hero = nil
	g.Deck = nil
	g.DeckInBoard = nil
	g.DeckInGraveYard = nil
}

// SetHero ...初始化群组中的英雄卡牌
func (g *Group) SetHero(c *Card) bool {
	if c.ActorType == 1 {
		gb := &GameBattler{GameBattlerBase: NewGameBattlerBase(), GameAction: new(GameAction)}
		gb.CopyCard(c)
		g.Hero = gb
		return true
	}
	return false
}

// SetDecks ...初始化群组中的卡牌
func (g *Group) SetDecks(d *Deck) bool {
	for _, v := range d.Cards {
		if v.ActorType == 2 {
			gb := &GameBattler{
				GameBattlerBase: NewGameBattlerBase(),
				GameAction:      &GameAction{},
			}
			gb.CopyCard(v)
			g.Deck = append(g.Deck, gb)
		} else {
			g.Deck = nil
			return false
		}
	}
	return true
}

// SummonCard ...打出手牌中冷却时间已经结束的卡牌
func (g *Group) SummonCard() {
	tmpHand := []*GameBattler{}
	tmpBoard := g.DeckInBoard
	for _, v := range g.Deck {
		if v.initTurnCooldown == 0 {
			v.IsInBoard = true
			v.Position = len(tmpBoard)
			tmpBoard = append(tmpBoard, v)
		} else {
			tmpHand = append(tmpHand, v)
		}
	}
	g.Deck = tmpHand
	g.DeckInBoard = tmpBoard
	//触发事件，动画等
}

// RemoveCard ...卡牌死亡移入墓地
func (g *Group) RemoveCard(c *GameBattler) {
	for k, v := range g.DeckInBoard {
		if v == c {
			c.IsInGraveYard = true
			g.DeckInGraveYard = append(g.DeckInGraveYard, c)
			if k == len(g.DeckInBoard)-1 {
				g.DeckInBoard = g.DeckInBoard[:k]
			} else {
				g.DeckInBoard = append(g.DeckInBoard[:k], g.DeckInBoard[k+1:]...)
			}
			break
		}
	}

	for k, v := range g.DeckInBoard {
		v.Position = k
	}
}

// CheckDeadCardsAndRemoveThem ...检查死亡的卡牌，并且移动到墓地中
func (g *Group) CheckDeadCardsAndRemoveThem() {
	tmpCards := []*GameBattler{}
	for _, v := range g.DeckInBoard {
		if v.IsDead() {
			tmpCards = append(tmpCards, v)
			fmt.Print("\nUser")
			fmt.Print(g.UserID)
			fmt.Print(" " + v.Name + " dead\n")
		}
	}
	if len(tmpCards) > 0 {
		for _, v := range tmpCards {
			g.RemoveCard(v)
		}
	}
}

func (g *Group) MinusTurnCoolDown() {
	for _, v := range g.Deck {
		v.MinusTurnCoolDown()
	}
}
