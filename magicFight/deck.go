package main

type Deck struct {
	ID            string
	Name          string
	MaxCardNumber int
	Card          []*Card
}

func (d *Deck) AddCard(c ...*Card) {
	for _, v := range c {
		if !d.ContainCard(v) {
			d.Card = append(d.Card, v)
		}
	}
}

func (d *Deck) ContainCard(c *Card) bool {
	for _, v := range d.Card {
		if v == c {
			return true
		}
	}
	return false
}

func (d *Deck) IsEmpty() bool {
	if len(d.Card) == 0 {
		return true
	}
	return false
}
