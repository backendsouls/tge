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
	STR float64 // Strength
	AGI float64 // Agility
	INT float64 // Intelligence
	VIT float64 // Vitality
	DEX float64 // Dexterity
	WIS float64 // Wisdom
	CHA float64 // Charisma
	LUK float64 // Luck
}

// NewStats validates and builds a stat block; every attribute must be >= 0.
func NewStats(str, agi, intel, vit, dex, wis, cha, luk float64) (Stats, error) {
	s := Stats{STR: str, AGI: agi, INT: intel, VIT: vit, DEX: dex, WIS: wis, CHA: cha, LUK: luk}
	for _, v := range []float64{str, agi, intel, vit, dex, wis, cha, luk} {
		if v < 0 {
			return Stats{}, ErrNegativeStat
		}
	}
	return s, nil
}

// BaseStats returns the default starting attributes for a new character.
func BaseStats() Stats {
	return Stats{STR: 0.65, AGI: 0.65, INT: 0.65, VIT: 0.65, DEX: 0.65, WIS: 0.65, CHA: 0.65, LUK: 0.65}
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
