package cultivation

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidLevelName is returned when a level name is blank.
	ErrInvalidLevelName = errors.New("level: name must not be empty")
	// ErrInvalidLevelNumber is returned when a level number is not positive.
	ErrInvalidLevelNumber = errors.New("level: number must be positive")
	// ErrInvalidLevelPoints is returned when a level's breakthrough points are negative.
	ErrInvalidLevelPoints = errors.New("level: points must not be negative")
	// ErrLevelNumberExists is returned when a level number is reused within a realm.
	ErrLevelNumberExists = errors.New("level: number already used in this realm")
)

// Level is a single stage within a Realm — e.g. the "First Level" of the "Qi
// Condensation" realm. Levels are ordered by Number within their realm. Unlike a
// Realm, a Level carries progression concepts: BreakthroughPoints to advance
// to the next level.
type Level struct {
	Number             int
	Name               string
	BreakthroughPoints int
}

// NewLevel validates and builds a level. The number must be positive, the name
// non-blank, and the points non-negative.
func NewLevel(number int, name string, breakthroughPoints int) (Level, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Level{}, ErrInvalidLevelName
	}
	if number <= 0 {
		return Level{}, fmt.Errorf("%w: %d", ErrInvalidLevelNumber, number)
	}
	if breakthroughPoints < 0 {
		return Level{}, ErrInvalidLevelPoints
	}
	return Level{
		Number:             number,
		Name:               name,
		BreakthroughPoints: breakthroughPoints,
	}, nil
}
