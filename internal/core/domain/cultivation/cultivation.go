package cultivation

import (
	"errors"
	"fmt"

	"tge/internal/core/domain/power"
	"tge/internal/core/domain/powersystem"
)

// ErrInvalidProgress is returned when cultivation progress is negative.
var ErrInvalidProgress = errors.New("cultivation: progress must not be negative")

// CultivationState is the Cultivation kind's PowerState: a character's standing at
// one power node, anchored to a Realm and the current Level within it, plus the
// points accumulated toward the next breakthrough. Power and lifespan are derived
// from the realm's ax+b formulas evaluated at Progress.
type CultivationState struct {
	Realm    Realm
	Level    Level   // current level within the realm
	Points   int     // accumulated toward Level.BreakthroughPoints (the breakthrough gate)
	Progress float64 // x fed into the realm's power/lifespan formulas
}

// NewCultivationState validates and builds a CultivationState. It returns
// ErrInvalidProgress if progress is negative.
func NewCultivationState(realm Realm, level Level, points int, progress float64) (CultivationState, error) {
	if progress < 0 {
		return CultivationState{}, fmt.Errorf("%w: %g", ErrInvalidProgress, progress)
	}
	return CultivationState{Realm: realm, Level: level, Points: points, Progress: progress}, nil
}

// Ready reports whether the breakthrough gate of the current level is full — the character
// has broken through this level and is ready to advance.
func (c CultivationState) Ready() bool {
	return c.Points >= c.Level.BreakthroughPoints
}

// AdvanceWithin accumulates points into the state, filling the current level's
// breakthrough gate. When the gate fills it breaks through to the next level
// in this realm (resetting the gate) and continues. It returns the updated state
// and any points left over after the realm's final level is fully broken through,
// for the caller to carry into the next realm.
func (c CultivationState) AdvanceWithin(points int) (CultivationState, int) {
	n := points
	for {
		switch {
		case c.Points < c.Level.BreakthroughPoints:
			if n == 0 {
				return c, 0
			}
			add := min(n, c.Level.BreakthroughPoints-c.Points)
			c.Points += add
			n -= add
		default: // gate full — break through to the next level in this realm
			next, ok := c.Realm.nextLevel(c.Level.Number)
			if !ok {
				return c, n // realm ceiling reached; hand leftover to the caller
			}
			c.Level = next
			c.Points = 0
		}
	}
}

// Kind implements PowerState.
func (CultivationState) Kind() powersystem.PowerSystemType { return powersystem.Cultivation }

// Power returns the cultivation's current power via the realm's formula.
func (c CultivationState) Power() float64 { return c.Realm.Power(c.Progress) }

// Lifespan returns the cultivation's current lifespan via the realm's formula.
func (c CultivationState) Lifespan() float64 { return c.Realm.Lifespan(c.Progress) }

var _ power.PowerState = CultivationState{}
