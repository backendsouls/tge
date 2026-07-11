package character

import (
	"errors"
	"fmt"
	"strings"
	"tge/internal/core/domain/power"
	"tge/internal/core/domain/rpg"
)

var (
	ErrInvalidName           = errors.New("character: name must not be empty")
	ErrInvalidCharacterType  = errors.New("character: invalid type")
	ErrInvalidGender         = errors.New("character: invalid gender")
	ErrMainCharacterRequired = errors.New("character: role requires a main character")
	ErrRoleGenderMismatch    = errors.New("character: role not allowed for the main character's gender")
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

func (g Gender) Valid() bool { return g == Male || g == Female }

type CharacterType string

const (
	MainCharacter    CharacterType = "MainCharacter"
	SideCharacter    CharacterType = "SideCharacter"
	SupportCharacter CharacterType = "SupportCharacter"
	Hero             CharacterType = "Hero"
	Heroine          CharacterType = "Heroine"
)

func (t CharacterType) Valid() bool {
	switch t {
	case MainCharacter, SideCharacter, SupportCharacter, Hero, Heroine:
		return true
	}
	return false
}

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

type Mortal struct {
	Age      int
	Lifespan int
}

// NodeProgress tracks a character's state at a specific PowerNode in a PowerSystem.
type NodeProgress struct {
	System    string
	NodeID    string
	Level     int
	Progress  float64
	BasePower float64
}

type IdleRates struct {
	SkillPointsPerHour       float64 `json:"skill_points_per_hour"`
	CultivationPointsPerHour float64 `json:"cultivation_points_per_hour"`
	AbilityPointsPerHour     float64 `json:"ability_points_per_hour"`
	ProfessionPointsPerHour  float64 `json:"profession_points_per_hour"`
}

type IdleState struct {
	StartTime      int64     `json:"start_time"`
	ActiveActivity string    `json:"active_activity"`
	Rates          IdleRates `json:"rates"`
}

type Character struct {
	Name          string
	Type          CharacterType
	Gender        Gender
	Species       []Species
	PowerValue    string // The string representation of Total Power
	MechanicState power.MechanicState
	UnlockedNodes []NodeProgress
	Mortal        Mortal
	Class         rpg.Class
	Profession    rpg.Profession
	Stats         rpg.Stats
	Inventory     rpg.Inventory
	IdleState     IdleState
	NovelTime     int64 // The current time in the novel/history for this character (in seconds)
}

// CalculateTotalPower dynamically computes the character's power level.
// Formula: (MechanicState.BasePower + Sum(UnlockedNodes.BasePower)) * Species.Power
func (c *Character) CalculateTotalPower() float64 {
	total := c.MechanicState.BasePower

	for _, node := range c.UnlockedNodes {
		total += node.BasePower
	}

	// Apply species multiplier
	if len(c.Species) > 0 {
		total *= c.Species[0].Power
	}

	if total < 1.0 {
		total = 1.0
	}

	c.PowerValue = fmt.Sprintf("%g", total)
	return total
}

// CurrentEnergyPools calculates the current energy pools including dynamically generated idle points.
func (c *Character) CurrentEnergyPools(currentNovelTime int64) map[string]int {
	pools := make(map[string]int)
	if c.MechanicState.EnergyPools != nil {
		for k, v := range c.MechanicState.EnergyPools {
			pools[k] = v
		}
	}

	if c.IdleState.StartTime >= 0 && currentNovelTime > c.IdleState.StartTime {
		deltaHours := float64(currentNovelTime-c.IdleState.StartTime) / 3600.0

		if c.IdleState.Rates.SkillPointsPerHour > 0 {
			pools["SkillPoints"] += int(deltaHours * c.IdleState.Rates.SkillPointsPerHour)
		}
		if c.IdleState.Rates.CultivationPointsPerHour > 0 {
			pools["CultivationPoints"] += int(deltaHours * c.IdleState.Rates.CultivationPointsPerHour)
		}
		if c.IdleState.Rates.AbilityPointsPerHour > 0 {
			pools["AbilityPoints"] += int(deltaHours * c.IdleState.Rates.AbilityPointsPerHour)
		}
		if c.IdleState.Rates.ProfessionPointsPerHour > 0 {
			pools["ProfessionPoints"] += int(deltaHours * c.IdleState.Rates.ProfessionPointsPerHour)
		}
	}

	return pools
}

type CharacterConfig struct {
	Name       string
	Type       CharacterType
	Gender     Gender
	Species    Species
	Class      rpg.Class
	Profession rpg.Profession
	Age        int
	Stats      rpg.Stats
}

const defaultAge = 16

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

	age := cfg.Age
	if age <= 0 {
		age = defaultAge
	}
	stats := cfg.Stats
	if stats == (rpg.Stats{}) {
		stats = rpg.BaseStats()
	}

	state, _ := power.NewMechanicState(0, 1.0)

	char := Character{
		Name:          name,
		Type:          cfg.Type,
		Gender:        cfg.Gender,
		Species:       []Species{cfg.Species},
		MechanicState: state,
		UnlockedNodes: []NodeProgress{},
		Mortal:        Mortal{Age: age, Lifespan: cfg.Species.Lifespan},
		Class:         cfg.Class,
		Profession:    cfg.Profession,
		Stats:         stats,
	}
	char.CalculateTotalPower()

	return char, nil
}

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
