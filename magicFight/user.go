package main

type User struct {
	ID          int
	Name        string
	SessionID   int
	DefaultDeck *Deck
	DefaultHero *Card
}

func (u *User) MakeOneDeck() {
	d := &Deck{}
	for _, v := range cardSource {
		if v.ActorType != 1 {
			t := v
			d.AddCard(&t)
		}
	}
	u.DefaultDeck = d
}

func (u *User) MakeOneHero() {
	for _, v := range cardSource {
		if v.ActorType == 1 {
			if v.Name == "Athena" {
				t := v
				u.DefaultHero = &t
			}
		}
	}
}
