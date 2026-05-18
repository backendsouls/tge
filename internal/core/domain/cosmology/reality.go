package cosmology

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidRealityName     = errors.New("reality: name must not be empty")
	ErrRealityOmniverseExists = errors.New("reality: omniverse already in this reality")
)

// Reality, also known as the Box, is a collection of Omniverses.
type Reality struct {
	Name       string
	Omniverses []Omniverse
}

// NewReality validates and creates a new Reality.
func NewReality(name string) (Reality, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Reality{}, ErrInvalidRealityName
	}
	return Reality{Name: name}, nil
}

// AddOmniverse adds an omniverse (by name) to the reality, rejecting duplicates.
func (r *Reality) AddOmniverse(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidOmniverseName
	}
	for _, o := range r.Omniverses {
		if o.Name == name {
			return fmt.Errorf("%w: %q", ErrRealityOmniverseExists, name)
		}
	}
	r.Omniverses = append(r.Omniverses, Omniverse{Name: name})
	return nil
}
