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

// PowerSystemType identifies the family a power system belongs to. Different kinds
// progress by different rules — a Cultivation system advances through realms and
// levels, a Magic system will advance by its own mechanics — so a character's
// per-node progress (PowerState) is kind-specific while the surrounding
// PowerSystem stays general.
type PowerSystemType string

const (
	// Cultivation is the realm/level progression kind.
	Cultivation PowerSystemType = "Cultivation"
	// Magic is a placeholder for a future, differently-progressing kind.
	Magic PowerSystemType = "Magic"
	// SuperPower is the new power system type.
	SuperPower PowerSystemType = "SuperPower"
)

// Valid reports whether k is a known system kind.
func (k PowerSystemType) Valid() bool {
	switch k {
	case Cultivation, Magic, SuperPower:
		return true
	}
	return false
}

// PowerState is a character's progress at a single power node. Each PowerSystemType
// supplies its own implementation (CultivationState today, a Magic state later),
// which keeps a character's attained power (Character.Power) general over the
// kind of system rather than tied to cultivation.
type PowerState interface {
	// Kind reports which system kind this state belongs to.
	Kind() PowerSystemType
	// Power reports the numerical progress of the state.
	Power() float64
}

// PowerSystem is a named tree (forest) of powers. For example "cosmology.Universe A
// Cultivation" might contain only a "Body" power, while another system contains
// BodyCultivation, SoulCultivation and SpiritCultivation, each with their own sub-powers.
type PowerSystem struct {
	Name            string
	PowerSystemType PowerSystemType // the family this system belongs to (defaults to Cultivation)
	Powers          []Power         // root powers
}

// NewPowerSystem validates and builds an empty power system.
func NewPowerSystem(name string, kind PowerSystemType) (PowerSystem, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return PowerSystem{}, ErrInvalidSystemName
	}
	if string(kind) == "" {
		kind = Cultivation
	}
	if !kind.Valid() {
		return PowerSystem{}, fmt.Errorf("power system: invalid kind %q", kind)
	}
	return PowerSystem{Name: name, PowerSystemType: kind}, nil
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
