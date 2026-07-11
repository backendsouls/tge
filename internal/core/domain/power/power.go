package power

import (
	"errors"
	"strings"

	"tge/internal/core/domain/powersystem"
)

// ErrInvalidPowerName is returned when a power is created with a blank name.
var ErrInvalidPowerName = errors.New("power: name must not be empty")

type PowerState interface {
	Kind() powersystem.PowerSystemType
	Power() float64
}

// Power is a node in a power system's tree. A power may contain child powers,
// letting a system model arbitrarily nested abilities or paths (e.g. a "Body"
// power with specific techniques beneath it). Powers are data created at runtime
// rather than Go types, so new powers require no code changes (OCP).
type Power struct {
	Name     string
	Children []Power
	// State is the owning character's progress at this node. It is nil in a power
	// system's definition and set only in a character's instance (Character.Power).
	State PowerState
}

// NewPower validates and builds a leaf power.
func NewPower(name string) (Power, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Power{}, ErrInvalidPowerName
	}
	return Power{Name: name}, nil
}
