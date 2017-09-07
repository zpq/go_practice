package main

type Deck struct {
	ID            string
	Name          string
	MaxCardNumber int
	Cards         []*Card
}

func (d *Deck) AddCard(c *Card) {
	if !d.ContainCard(c) {
		d.Cards = append(d.Cards, c)
	}
}

func (d *Deck) ContainCard(c *Card) bool {
	for _, v := range d.Cards {
		if v == c {
			return true
		}
	}
	return false
}

func (d *Deck) RemoveCard(c *Card) {

}

func (d *Deck) IsEmpty() bool {
	if len(d.Cards) == 0 {
		return true
	}
	return false
}
