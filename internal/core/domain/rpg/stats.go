// Package rpg holds the role-playing-game building blocks: character-defining
// entities (Class, Profession, Stats), things a character can have or do
// (Ability, Skill, Item, Equipment, Inventory), modifiers (Effect) and
// content (Quest, Recipe). Entities are keyed by Name, mirroring the rest of
// the domain.
package rpg

import "errors"

// ErrNegativeStat is returned when a stat attribute is negative.
var ErrNegativeStat = errors.New("stats: attribute must not be negative")

// Stats holds a character's eight core attributes.
type Stats struct {
	STR int // Strength
	AGI int // Agility
	INT int // Intelligence
	VIT int // Vitality
	DEX int // Dexterity
	WIS int // Wisdom
	CHA int // Charisma
	LUK int // Luck
}

// NewStats validates and builds a stat block; every attribute must be >= 0.
func NewStats(str, agi, intel, vit, dex, wis, cha, luk int) (Stats, error) {
	s := Stats{STR: str, AGI: agi, INT: intel, VIT: vit, DEX: dex, WIS: wis, CHA: cha, LUK: luk}
	for _, v := range []int{str, agi, intel, vit, dex, wis, cha, luk} {
		if v < 0 {
			return Stats{}, ErrNegativeStat
		}
	}
	return s, nil
}

// BaseStats returns the default starting attributes for a new character.
func BaseStats() Stats {
	return Stats{STR: 5, AGI: 5, INT: 5, VIT: 5, DEX: 5, WIS: 5, CHA: 5, LUK: 5}
}

// Add returns the element-wise sum of two stat blocks (e.g. base + equipment
// bonus).
func (s Stats) Add(o Stats) Stats {
	return Stats{
		STR: s.STR + o.STR,
		AGI: s.AGI + o.AGI,
		INT: s.INT + o.INT,
		VIT: s.VIT + o.VIT,
		DEX: s.DEX + o.DEX,
		WIS: s.WIS + o.WIS,
		CHA: s.CHA + o.CHA,
		LUK: s.LUK + o.LUK,
	}
}
