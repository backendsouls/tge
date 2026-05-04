// Package domain holds the cultivation game's core entities and business rules.
//
// It is the center of the hexagon: it depends on nothing outside itself, so
// every other layer (ports and adapters) depends inward on it.
package progression

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	// ErrInvalidName is returned when a realm name is empty or whitespace only.
	ErrInvalidName = errors.New("realm: name must not be empty")
	// ErrInvalidPoints is returned when a realm's bottleneck points are negative.
	ErrInvalidPoints = errors.New("realm: points must not be negative")
	// ErrInvalidMaxLevels is returned when a realm's maximum level count is negative.
	ErrInvalidMaxLevels = errors.New("realm: max levels must not be negative")
	// ErrLevelNumberExceedsMax is returned when a level number is above the realm's MaxLevels cap.
	ErrLevelNumberExceedsMax = errors.New("level: number exceeds the realm's maximum levels")
)

// Realm represents a single stage of cultivation.
//
// Power and Lifespan each follow a linear ax + b formula: the multiplier is the
// slope (a) and the adder is the intercept (b), evaluated against a cultivation
// progress value x.
type Realm struct {
	Name string

	PowerMultiplier float64 // a in the power formula
	PowerAdder      float64 // b in the power formula

	LifespanMultiplier float64 // a in the lifespan formula
	LifespanAdder      float64 // b in the lifespan formula

	// BottleneckPoints is the realm-wide difficulty barrier. Breakthrough is a
	// per-level concept (see Level), not a realm one.
	BottleneckPoints int

	// MaxLevels is how many levels a normal character may reach in this realm
	// (0 = unlimited). MainCharacterMaxLevels is the main character's higher cap;
	// 0 means the main character uses MaxLevels (no special privilege).
	MaxLevels              int
	MainCharacterMaxLevels int
	Levels                 []Level // ordered sub-stages within this realm
}

// MaxLevelsFor returns the level cap for a character, given whether it is the
// main character. A main character uses MainCharacterMaxLevels when set,
// otherwise the normal MaxLevels. A returned 0 means unlimited.
func (r Realm) MaxLevelsFor(isMain bool) int {
	if isMain && r.MainCharacterMaxLevels > 0 {
		return r.MainCharacterMaxLevels
	}
	return r.MaxLevels
}

// effectiveMaxLevels is the most levels this realm may define — enough for the
// most-privileged character (the main character) to reach. 0 means unlimited.
func (r Realm) effectiveMaxLevels() int {
	if r.MaxLevels == 0 {
		return 0 // normal tier unlimited ⇒ the realm is unlimited
	}
	if r.MainCharacterMaxLevels > r.MaxLevels {
		return r.MainCharacterMaxLevels
	}
	return r.MaxLevels
}

// AddLevel adds an ordered level to the realm, rejecting a blank name, a
// non-positive number, negative points, a duplicate number, or a number above
// the realm's effective cap (the higher, main-character, cap when one is set).
// Levels are kept sorted by Number.
func (r *Realm) AddLevel(number int, name string, breakthroughPoints, bottleneckPoints int) error {
	lvl, err := NewLevel(number, name, breakthroughPoints, bottleneckPoints)
	if err != nil {
		return err
	}
	if cap := r.effectiveMaxLevels(); cap > 0 && lvl.Number > cap {
		return fmt.Errorf("%w: %d (max %d)", ErrLevelNumberExceedsMax, lvl.Number, cap)
	}
	for _, l := range r.Levels {
		if l.Number == lvl.Number {
			return fmt.Errorf("%w: %d", ErrLevelNumberExists, lvl.Number)
		}
	}
	r.Levels = append(r.Levels, lvl)
	sort.Slice(r.Levels, func(i, j int) bool { return r.Levels[i].Number < r.Levels[j].Number })
	return nil
}

// RealmConfig holds the attributes used to construct a Realm.
type RealmConfig struct {
	Name                   string
	PowerMultiplier        float64
	PowerAdder             float64
	LifespanMultiplier     float64
	LifespanAdder          float64
	BottleneckPoints       int
	MaxLevels              int
	MainCharacterMaxLevels int
}

// NewRealm validates cfg and returns a Realm. It returns character.ErrInvalidName if the
// name is blank and ErrInvalidPoints if either points value is negative.
func NewRealm(cfg RealmConfig) (Realm, error) {
	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		return Realm{}, ErrInvalidName
	}
	if cfg.BottleneckPoints < 0 {
		return Realm{}, fmt.Errorf("%w: bottleneck points %d", ErrInvalidPoints, cfg.BottleneckPoints)
	}
	if cfg.MaxLevels < 0 {
		return Realm{}, fmt.Errorf("%w: %d", ErrInvalidMaxLevels, cfg.MaxLevels)
	}
	if cfg.MainCharacterMaxLevels < 0 {
		return Realm{}, fmt.Errorf("%w: main %d", ErrInvalidMaxLevels, cfg.MainCharacterMaxLevels)
	}

	return Realm{
		Name:                   name,
		PowerMultiplier:        cfg.PowerMultiplier,
		PowerAdder:             cfg.PowerAdder,
		LifespanMultiplier:     cfg.LifespanMultiplier,
		LifespanAdder:          cfg.LifespanAdder,
		BottleneckPoints:       cfg.BottleneckPoints,
		MaxLevels:              cfg.MaxLevels,
		MainCharacterMaxLevels: cfg.MainCharacterMaxLevels,
	}, nil
}

// Power returns the cultivator's power at progress x using ax + b.
func (r Realm) Power(x float64) float64 {
	return r.PowerMultiplier*x + r.PowerAdder
}

// Lifespan returns the cultivator's lifespan at progress x using ax + b.
func (r Realm) Lifespan(x float64) float64 {
	return r.LifespanMultiplier*x + r.LifespanAdder
}
