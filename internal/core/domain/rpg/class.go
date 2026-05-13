package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidClassName is returned when a class name is blank.
var ErrInvalidClassName = errors.New("class: name must not be empty")

// Class is a character's combat archetype (e.g. Warrior, Mage).
type Class struct {
	Name        string
	Description string
}

// NewClass validates and builds a class.
func NewClass(name, description string) (Class, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Class{}, ErrInvalidClassName
	}
	return Class{Name: name, Description: strings.TrimSpace(description)}, nil
}
