package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidClassName is returned when a class name is blank.
var ErrInvalidClassName = errors.New("class: name must not be empty")

// Class is an archetype that determines stat growth or capabilities.
type Class struct {
	Name        string
	Description string
	Grade       string
}

// NewClass validates and builds a class.
func NewClass(name, description, grade string) (Class, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Class{}, ErrInvalidClassName
	}
	return Class{Name: name, Description: strings.TrimSpace(description), Grade: strings.TrimSpace(grade)}, nil
}
