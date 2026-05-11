package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidSkillName is returned when a skill name is blank.
var ErrInvalidSkillName = errors.New("skill: name must not be empty")

// Skill is a learned, trainable ability.
type Skill struct {
	Name        string
	Description string
}

// NewSkill validates and builds a skill.
func NewSkill(name, description string) (Skill, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Skill{}, ErrInvalidSkillName
	}
	return Skill{Name: name, Description: strings.TrimSpace(description)}, nil
}
