package progression

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidSystemName is returned when a power system has a blank name.
	ErrInvalidSystemName = errors.New("power system: name must not be empty")
	// ErrPowerExists is returned when adding a power whose name is already used in the system.
	ErrPowerExists = errors.New("power system: power already exists")
	// ErrPowerParentNotFound is returned when the named parent power does not exist.
	ErrPowerParentNotFound = errors.New("power system: parent power not found")
)

// PowerSystem is a named tree (forest) of powers. For example "cosmology.Universe A
// Cultivation" might contain only a "Body" power, while another system contains
// Body, Soul and Spirit, each with their own sub-powers.
type PowerSystem struct {
	Name   string
	Kind   SystemKind // the family this system belongs to (defaults to Cultivation)
	Powers []Power    // root powers
}

// NewPowerSystem validates and builds an empty power system of the Cultivation
// kind, the default until other kinds (e.g. Magic) are authored.
func NewPowerSystem(name string) (PowerSystem, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return PowerSystem{}, ErrInvalidSystemName
	}
	return PowerSystem{Name: name, Kind: Cultivation}, nil
}

// AddPower inserts a power into the system. With an empty parent the power
// becomes a root; otherwise it is attached beneath the named parent (found
// anywhere in the tree). Power names must be unique across the whole system.
func (ps *PowerSystem) AddPower(name, parent string) error {
	power, err := NewPower(name)
	if err != nil {
		return err
	}
	if _, ok := findInForest(ps.Powers, power.Name); ok {
		return fmt.Errorf("%w: %q", ErrPowerExists, power.Name)
	}
	if strings.TrimSpace(parent) == "" {
		ps.Powers = append(ps.Powers, power)
		return nil
	}
	p, ok := findPtrInForest(ps.Powers, parent)
	if !ok {
		return fmt.Errorf("%w: %q", ErrPowerParentNotFound, parent)
	}
	p.Children = append(p.Children, power)
	return nil
}

// Names returns every power name in pre-order, useful for display and tests.
func (ps PowerSystem) Names() []string {
	var out []string
	var walk func([]Power)
	walk = func(powers []Power) {
		for _, p := range powers {
			out = append(out, p.Name)
			walk(p.Children)
		}
	}
	walk(ps.Powers)
	return out
}

func findInForest(powers []Power, name string) (Power, bool) {
	for _, p := range powers {
		if p.Name == name {
			return p, true
		}
		if c, ok := findInForest(p.Children, name); ok {
			return c, true
		}
	}
	return Power{}, false
}

func findPtrInForest(powers []Power, name string) (*Power, bool) {
	for i := range powers {
		if powers[i].Name == name {
			return &powers[i], true
		}
		if c, ok := findPtrInForest(powers[i].Children, name); ok {
			return c, true
		}
	}
	return nil, false
}
