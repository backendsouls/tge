package superpower

import (
	"errors"

	"tge/internal/core/domain/powersystem"
)

var ErrInvalidTier = errors.New("super power: tier must be between 0 and 5")

type SuperPowerState struct {
	Tier int // 0 to 5
}

// NewSuperPowerState creates a valid SuperPowerState, validating the tier bounds.
func NewSuperPowerState(tier int) (SuperPowerState, error) {
	if tier < 0 || tier > 5 {
		return SuperPowerState{}, ErrInvalidTier
	}
	return SuperPowerState{Tier: tier}, nil
}

func (SuperPowerState) Kind() powersystem.PowerSystemType { return powersystem.SuperPower }

func (s SuperPowerState) Power() float64 {
	switch s.Tier {
	case 0:
		return 1.0 // Base
	case 1:
		return 5.0
	case 2:
		return 10.0
	case 3:
		return 20.0
	case 4:
		return 50.0
	case 5:
		return 100.0
	}
	return 1.0
}
