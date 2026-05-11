package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidAbilityName is returned when an ability name is blank.
var ErrInvalidAbilityName = errors.New("ability: name must not be empty")

// Ability is an innate power a character possesses.
type Ability struct {
	Name        string
	Description string
}

// NewAbility validates and builds an ability.
func NewAbility(name, description string) (Ability, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Ability{}, ErrInvalidAbilityName
	}
	return Ability{Name: name, Description: strings.TrimSpace(description)}, nil
}
