package rpg

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidEffectName is returned when an effect name is blank.
	ErrInvalidEffectName = errors.New("effect: name must not be empty")
	// ErrInvalidEffectKind is returned for an unknown effect kind.
	ErrInvalidEffectKind = errors.New("effect: invalid kind")
)

// EffectKind classifies an effect.
type EffectKind string

const (
	Buff   EffectKind = "Buff"   // a positive modifier
	Debuff EffectKind = "Debuff" // a negative modifier
	Status EffectKind = "Status" // a state (e.g. Poisoned, Stunned)
)

// Valid reports whether k is a known effect kind.
func (k EffectKind) Valid() bool {
	switch k {
	case Buff, Debuff, Status:
		return true
	}
	return false
}

// Effect is a modifier or condition that abilities, skills or items can apply.
type Effect struct {
	Name        string
	Kind        EffectKind
	Description string
}

// NewEffect validates and builds an effect.
func NewEffect(name string, kind EffectKind, description string) (Effect, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Effect{}, ErrInvalidEffectName
	}
	if !kind.Valid() {
		return Effect{}, fmt.Errorf("%w: %q", ErrInvalidEffectKind, kind)
	}
	return Effect{Name: name, Kind: kind, Description: strings.TrimSpace(description)}, nil
}
