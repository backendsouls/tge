package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidProfessionName is returned when a profession name is blank.
var ErrInvalidProfessionName = errors.New("profession: name must not be empty")

// Profession is a non-combat or specialized career path.
type Profession struct {
	Name        string
	Description string
	Grade       string
}

// NewProfession validates and builds a profession.
func NewProfession(name, description, grade string) (Profession, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Profession{}, ErrInvalidProfessionName
	}
	return Profession{Name: name, Description: strings.TrimSpace(description), Grade: strings.TrimSpace(grade)}, nil
}
