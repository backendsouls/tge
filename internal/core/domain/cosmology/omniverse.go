package cosmology

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidOmniverseName      = errors.New("omniverse: name must not be empty")
	ErrOmniverseMultiverseExists = errors.New("omniverse: multiverse already in this omniverse")
)

// Omniverse is a collection of Multiverses.
type Omniverse struct {
	Name        string
	Multiverses []Multiverse
}

// NewOmniverse validates and creates a new Omniverse.
func NewOmniverse(name string) (Omniverse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Omniverse{}, ErrInvalidOmniverseName
	}
	return Omniverse{Name: name}, nil
}

// AddMultiverse adds a multiverse (by name) to the omniverse, rejecting duplicates.
func (o *Omniverse) AddMultiverse(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidMultiverseName
	}
	for _, m := range o.Multiverses {
		if m.Name == name {
			return fmt.Errorf("%w: %q", ErrOmniverseMultiverseExists, name)
		}
	}
	o.Multiverses = append(o.Multiverses, Multiverse{Name: name})
	return nil
}
