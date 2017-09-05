package main

// Creature ...base model
type Creature struct {
	ID      string
	Name    string
	Level   int
	Hp      int
	Attack  int
	Rarity  int
	Faction int
	Star    int
	Skill   []*Skill
	Buff    []*Buff
}

func (c *Creature) GetFields() *Creature {
	return c
}
