package cosmology

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidMultiverseName    = errors.New("multiverse: name must not be empty")
	ErrMultiverseUniverseExists = errors.New("multiverse: universe already in this multiverse")
)

// Multiverse is a collection of Universes.
type Multiverse struct {
	Name      string
	Universes []Universe
}

// NewMultiverse validates and creates a new Multiverse.
func NewMultiverse(name string) (Multiverse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Multiverse{}, ErrInvalidMultiverseName
	}
	return Multiverse{Name: name}, nil
}

// AddUniverse adds a universe (by name) to the multiverse, rejecting duplicates.
func (m *Multiverse) AddUniverse(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidUniverseName
	}
	for _, u := range m.Universes {
		if u.Name == name {
			return fmt.Errorf("%w: %q", ErrMultiverseUniverseExists, name)
		}
	}
	m.Universes = append(m.Universes, Universe{Name: name})
	return nil
}
