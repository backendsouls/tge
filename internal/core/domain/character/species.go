package character

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSpeciesName = errors.New("species: name must not be empty")
	// ErrInvalidSpeciesGender is returned when a species' default gender is not a known gender.
	ErrInvalidSpeciesGender = errors.New("species: invalid default gender")
)

// HumanBaseName is the name of the built-in Human base species.
const HumanBaseName = "Human"

// HumanBase returns the built-in Human base species. It is the template for a
// newly created main character's initial mortal status (base Power and Lifespan)
// and defaults its members to Male.
func HumanBase() Species {
	return Species{Name: HumanBaseName, Power: 0.65, Lifespan: 80, DefaultGender: Male}
}

// Species represents a biological classification with base status values and a
// default gender used when a character of this species is created without one.
type Species struct {
	Name          string
	Power         float64
	Lifespan      int
	DefaultGender Gender // optional; empty means no species-level default
}

// NewSpecies validates and creates a new Species. defaultGender may be empty
// (no default) but must otherwise be a known gender.
func NewSpecies(name string, power float64, lifespan int, defaultGender Gender) (Species, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Species{}, ErrInvalidSpeciesName
	}
	if defaultGender != "" && !defaultGender.Valid() {
		return Species{}, fmt.Errorf("%w: %q", ErrInvalidSpeciesGender, defaultGender)
	}
	return Species{Name: name, Power: power, Lifespan: lifespan, DefaultGender: defaultGender}, nil
}
