package cosmology

import (
	"errors"
	"fmt"
	"strings"

	"tge/internal/core/domain/powersystem"
)

var (
	// ErrInvalidUniverseName is returned when a universe name is blank.
	ErrInvalidUniverseName = errors.New("universe: name must not be empty")
	// ErrUniverseSystemExists is returned when a power system is already in the universe.
	ErrUniverseSystemExists = errors.New("universe: power system already in this universe")
	// ErrRealmExistsInUniverse is returned when a realm (location) name is already used in the universe.
	ErrRealmExistsInUniverse = errors.New("universe: realm already exists in this universe")
	// ErrSingleRealm is returned when a bubble universe is initialized with more than one realm.
	ErrSingleRealm = errors.New("universe: bubble universe can only have one realm")
)

// Universe groups several power systems and contains in-universe realms
// (locations). Member systems are referenced by name; their power trees live
// with the power systems themselves.
type Universe struct {
	Name    string
	Systems []powersystem.PowerSystem
	Realms  []Location
}

// NewUniverse validates and builds an empty universe.
func NewUniverse(name string) (Universe, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Universe{}, ErrInvalidUniverseName
	}
	return Universe{Name: name}, nil
}

// AddSystem adds a power system (by name) to the universe, rejecting a duplicate
// within this universe.
func (u *Universe) AddSystem(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return powersystem.ErrInvalidSystemName
	}
	for _, s := range u.Systems {
		if s.Name == name {
			return fmt.Errorf("%w: %q", ErrUniverseSystemExists, name)
		}
	}
	u.Systems = append(u.Systems, powersystem.PowerSystem{Name: name})
	return nil
}

// AddRealms adds one or more in-universe realms (locations), rejecting
// duplicates. A universe may have any number of realms: zero, one (a "bubble"
// realm inside the universe), or many.
func (u *Universe) AddRealms(names ...string) error {
	seen := map[string]bool{}
	for _, r := range u.Realms {
		seen[r.Name] = true
	}
	additions := make([]Location, 0, len(names))
	for _, name := range names {
		loc, err := NewLocation(name)
		if err != nil {
			return err
		}
		if seen[loc.Name] {
			return fmt.Errorf("%w: %q", ErrRealmExistsInUniverse, loc.Name)
		}
		seen[loc.Name] = true
		additions = append(additions, loc)
	}
	u.Realms = append(u.Realms, additions...)
	return nil
}
