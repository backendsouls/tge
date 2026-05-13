package rpg

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidEquipmentName is returned when an equipment name is blank.
	ErrInvalidEquipmentName = errors.New("equipment: name must not be empty")
	// ErrInvalidEquipmentSlot is returned for an unknown equipment slot.
	ErrInvalidEquipmentSlot = errors.New("equipment: invalid slot")
)

// EquipmentSlot is where a piece of equipment is worn.
type EquipmentSlot string

const (
	Weapon    EquipmentSlot = "Weapon"
	Armor     EquipmentSlot = "Armor"
	Accessory EquipmentSlot = "Accessory"
)

// Valid reports whether s is a known equipment slot.
func (s EquipmentSlot) Valid() bool {
	switch s {
	case Weapon, Armor, Accessory:
		return true
	}
	return false
}

// Equipment is an equippable item occupying a slot and granting stat bonuses.
type Equipment struct {
	Name  string
	Slot  EquipmentSlot
	Bonus Stats
}

// NewEquipment validates and builds a piece of equipment.
func NewEquipment(name string, slot EquipmentSlot, bonus Stats) (Equipment, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Equipment{}, ErrInvalidEquipmentName
	}
	if !slot.Valid() {
		return Equipment{}, fmt.Errorf("%w: %q", ErrInvalidEquipmentSlot, slot)
	}
	return Equipment{Name: name, Slot: slot, Bonus: bonus}, nil
}
