package power

import (
	"errors"
)

var (
	ErrMechanicInvalidTier = errors.New("mechanic state: tier cannot be negative")
	ErrInvalidAlignment    = errors.New("mechanic state: alignment must be between -100 and 100")
)

// MechanicState replaces specific Cultivation/SuperPower states with a unified
// multi-dimensional progression tracker that sits at the character aggregate level.
type MechanicState struct {
	Tier            int
	BasePower       float64
	IsAwakened      bool
	Alignment       float64
	EnergyPools     map[string]int
	SpellSlots      map[int]int
	PermanentTraits []string
	Vows            []string
}

// NewMechanicState validates and initializes a character's universal progression state.
func NewMechanicState(tier int, basePower float64) (MechanicState, error) {
	if tier < 0 {
		return MechanicState{}, ErrMechanicInvalidTier
	}

	return MechanicState{
		Tier:            tier,
		BasePower:       basePower,
		IsAwakened:      false,
		Alignment:       0.0, // 0.0 is neutral
		EnergyPools:     make(map[string]int),
		SpellSlots:      make(map[int]int),
		PermanentTraits: []string{},
		Vows:            []string{},
	}, nil
}

// AddEnergyPool registers or updates a specific energy pool (e.g., Mana, Shinsu).
func (m *MechanicState) AddEnergyPool(name string, maxAmount int) {
	if m.EnergyPools == nil {
		m.EnergyPools = make(map[string]int)
	}
	m.EnergyPools[name] = maxAmount
}

// SetAlignment adjusts the character's alignment, strictly bounded between -100 and 100.
func (m *MechanicState) SetAlignment(val float64) error {
	if val < -100.0 || val > 100.0 {
		return ErrInvalidAlignment
	}
	m.Alignment = val
	return nil
}
