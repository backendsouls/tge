package cosmology

import (
	"errors"
	"strings"
)

// ErrInvalidLocationName is returned when a location (in-universe realm) name is blank.
var ErrInvalidLocationName = errors.New("location: name must not be empty")

// Location is a "realm" in the in-universe sense: a named place within a universe
// (e.g. "Hell progression.Realm", "character.Mortal progression.Realm", "Heaven progression.Realm").
//
// It is distinct from and unrelated to the cultivation progression.Realm, which is a power
// layer (ax+b). They merely share the word "realm".
type Location struct {
	Name string
}

// NewLocation validates and builds a location.
func NewLocation(name string) (Location, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Location{}, ErrInvalidLocationName
	}
	return Location{Name: name}, nil
}
