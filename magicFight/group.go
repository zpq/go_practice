package main

type Group struct {
	Hero            *GameBattler
	Deck            []*GameBattler
	DeckInBoard     []*GameBattler
	DeckInGraveYard []*GameBattler
	DeckInHand      []*GameBattler
}
