package character

import (
	"errors"
	"fmt"
	"math"
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
	System    string  `json:"system"`
	NodeID    string  `json:"node_id"`
	Level     int     `json:"level"`
	Progress  float64 `json:"progress"`
	BasePower float64 `json:"base_power"`
}

// CalculateBreakthrough returns the points required to breakthrough to the next level
func (n *NodeProgress) CalculateBreakthrough() float64 {
	return 100.0 * float64(n.Level*n.Level)
}

// Advance tries to consume points to level up. Returns any remaining (unconsumed) points.
func (n *NodeProgress) Advance(points float64) float64 {
	for points > 0 {
		reqBreakthrough := n.CalculateBreakthrough()
		if n.Progress < reqBreakthrough {
			add := math.Min(points, reqBreakthrough-n.Progress)
			n.Progress += add
			points -= add
			if points == 0 {
				return 0
			}
		}

		// Gate filled -> level up!
		n.Level++
		n.Progress = 0
	}
	return points
}

type IdleSlot struct {
	StartTime int64   `json:"start_time"`
	Duration  float64 `json:"duration"` // in hours, <= 0 means indefinite
	System    string  `json:"system"`
	Power     string  `json:"power"`
	Rate      float64 `json:"rate"`
}

type IdleState struct {
	Slots      []IdleSlot `json:"slots"`
	TotalSlots int        `json:"total_slots"`
}

type Character struct {
	Name          string
	Type          CharacterType
	Gender        Gender
	Species       []Species
	Systems       []string // Names of power systems this character belongs to
	PowerValue    string   // The string representation of Total Power
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
		total += node.BasePower * float64(node.Level)
	}

	for _, sp := range c.Species {
		total *= sp.Power
	}

	if total < 1.0 {
		total = 1.0
	}

	c.PowerValue = fmt.Sprintf("%.0f", total)
	return total
}

// AdvanceNode applies points to a specific NodeProgress. It returns any leftover points.
func (c *Character) AdvanceNode(system, nodeID string, points float64) float64 {
	for i, np := range c.UnlockedNodes {
		if np.System == system && np.NodeID == nodeID {
			remainder := np.Advance(points)
			c.UnlockedNodes[i] = np
			c.CalculateTotalPower()
			return remainder
		}
	}
	// Node not unlocked, all points remain
	return points
}

// CurrentEnergyPools calculates the current energy pools including dynamically generated idle points.
func (c *Character) CurrentEnergyPools(currentNovelTime int64) map[string]int {
	pools := make(map[string]int)
	if c.MechanicState.EnergyPools != nil {
		for k, v := range c.MechanicState.EnergyPools {
			pools[k] = v
		}
	}

	for _, slot := range c.IdleState.Slots {
		if slot.StartTime >= 0 && currentNovelTime > slot.StartTime {
			deltaHours := float64(currentNovelTime-slot.StartTime) / 3600.0
			if slot.Duration > 0 && deltaHours > slot.Duration {
				deltaHours = slot.Duration
			}
			if slot.Rate > 0 {
				poolName := fmt.Sprintf("%s_%s", slot.System, slot.Power)
				pools[poolName] += int(deltaHours * slot.Rate)
			}
		}
	}

	return pools
}

type CharacterConfig struct {
	Name       string
	Type       CharacterType
	Gender     Gender
	Species    Species
	Systems    []string
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

	mState, _ := power.NewMechanicState(0, 1.0)

	char := Character{
		Name:          name,
		Type:          cfg.Type,
		Gender:        cfg.Gender,
		Species:       []Species{cfg.Species},
		Systems:       cfg.Systems,
		MechanicState: mState,
		UnlockedNodes: []NodeProgress{},
		IdleState:     IdleState{TotalSlots: 1},
		Mortal: Mortal{
			Age:      age,
			Lifespan: cfg.Species.Lifespan,
		},
		Class:      cfg.Class,
		Profession: cfg.Profession,
		Stats:      stats,
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
