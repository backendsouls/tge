package character

import (
	"errors"
	"fmt"
	"strings"
	"tge/internal/core/domain/progression"
	"tge/internal/core/domain/rpg"
)

var (
	// ErrInvalidName is returned when a character name is empty or whitespace only.
	ErrInvalidName = errors.New("character: name must not be empty")
	// ErrMissingSystem is returned when a character has no power system.
	ErrMissingSystem = errors.New("character: at least one power system is required")
	// ErrInvalidCharacterType is returned for an unknown character type.
	ErrInvalidCharacterType = errors.New("character: invalid type")
	// ErrInvalidGender is returned for an unknown gender.
	ErrInvalidGender = errors.New("character: invalid gender")
	// ErrMainCharacterRequired is returned when a Hero/Heroine is created with no main character.
	ErrMainCharacterRequired = errors.New("character: role requires a main character")
	// ErrRoleGenderMismatch is returned when a role conflicts with the main character's gender.
	ErrRoleGenderMismatch = errors.New("character: role not allowed for the main character's gender")
)

// Gender is a character's gender. Hero/Heroine eligibility depends on the main
// character's gender.
type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

// Valid reports whether g is a known gender.
func (g Gender) Valid() bool { return g == Male || g == Female }

// CharacterType is a character's narrative role. It is a closed set (enum-like),
// unlike powers which are an open, growing tree of entities.
type CharacterType string

const (
	MainCharacter    CharacterType = "MainCharacter"
	SideCharacter    CharacterType = "SideCharacter"
	SupportCharacter CharacterType = "SupportCharacter"
	Hero             CharacterType = "Hero"    // male lead; requires a female main character
	Heroine          CharacterType = "Heroine" // female lead; requires a male main character
)

// Valid reports whether t is a known character type.
func (t CharacterType) Valid() bool {
	switch t {
	case MainCharacter, SideCharacter, SupportCharacter, Hero, Heroine:
		return true
	}
	return false
}

// requiredMainGender returns the main-character gender this role requires, and
// whether the role has such a requirement at all.
func (t CharacterType) requiredMainGender() (Gender, bool) {
	switch t {
	case Hero:
		return Female, true
	case Heroine:
		return Male, true
	default:
		return "", false
	}
}

// Mortal holds the non-cultivation attributes every character has from birth.
type Mortal struct {
	Age      int
	Lifespan int
}

// Character is a being in the story. It belongs to one or more power systems and
// is created in its initial mortal form (only mortal data, no cultivation). Its
// RPG attributes (Class, Profession, Stats, Inventory) are optional and default
// to empty/base values.
type Character struct {
	Name       string
	Type       CharacterType
	Gender     Gender
	Species    []Species
	PowerValue string // numeric value rendered as a string; computed later
	// Power is the character's attained power: its own progressed instances of the
	// power systems it has developed, one PowerSystem per system, with progress
	// hung on the tree nodes (PowerState). Empty for a fresh mortal.
	Power []progression.PowerSystem
	// PowerSystems is the complete tree definitions the character has access to.
	PowerSystems []progression.PowerSystem
	Mortal       Mortal
	Class        rpg.Class      // optional
	Profession   rpg.Profession // optional
	Stats        rpg.Stats
	Inventory    rpg.Inventory
}

// CharacterConfig holds the attributes used to construct a Character. Age and
// Stats are optional: a non-positive Age falls back to defaultAge and a zero
// Stats block falls back to rpg.BaseStats.
type CharacterConfig struct {
	Name         string
	Type         CharacterType
	Gender       Gender
	Species      Species
	PowerSystems []progression.PowerSystem
	Class        rpg.Class
	Profession   rpg.Profession
	Age          int
	Stats        rpg.Stats
}

// defaultAge is the starting age for a new mortal when none is configured.
const defaultAge = 16

// NewMortalCharacter validates cfg and builds a fresh mortal character. Cross-
// character rules (Hero/Heroine vs. the main character) are enforced separately
// via CheckRole, since they require knowledge of other characters.
func NewMortalCharacter(cfg CharacterConfig) (Character, error) {
	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		return Character{}, ErrInvalidName
	}
	if !cfg.Type.Valid() {
		return Character{}, fmt.Errorf("%w: %q", ErrInvalidCharacterType, cfg.Type)
	}
	if !cfg.Gender.Valid() {
		return Character{}, fmt.Errorf("%w: %q", ErrInvalidGender, cfg.Gender)
	}
	if len(cfg.PowerSystems) == 0 {
		return Character{}, ErrMissingSystem
	}
	power := fmt.Sprintf("%g", cfg.Species.Power)
	age := cfg.Age
	if age <= 0 {
		age = defaultAge
	}
	stats := cfg.Stats
	if stats == (rpg.Stats{}) {
		stats = rpg.BaseStats()
	}
	return Character{
		Name:         name,
		Type:         cfg.Type,
		Gender:       cfg.Gender,
		Species:      []Species{cfg.Species},
		PowerValue:   power,
		PowerSystems: cfg.PowerSystems,
		Mortal:       Mortal{Age: age, Lifespan: cfg.Species.Lifespan},
		Class:        cfg.Class,
		Profession:   cfg.Profession,
		Stats:        stats,
	}, nil
}

// CheckRole enforces the Hero/Heroine rules against the existing main characters.
// A lead role is allowed when at least one main character has the required
// gender. Non-lead roles always pass.
func CheckRole(t CharacterType, mainCharacters []Character) error {
	want, requires := t.requiredMainGender()
	if !requires {
		return nil
	}
	if len(mainCharacters) == 0 {
		return fmt.Errorf("%w: %s requires a main character to exist first", ErrMainCharacterRequired, t)
	}
	for _, mc := range mainCharacters {
		if mc.Gender == want {
			return nil
		}
	}
	return fmt.Errorf("%w: %s requires a %s main character", ErrRoleGenderMismatch, t, want)
}
