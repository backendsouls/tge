package progression

import (
	"errors"
	"fmt"
)

// ErrInvalidProgress is returned when cultivation progress is negative.
var ErrInvalidProgress = errors.New("cultivation: progress must not be negative")

// CultivationState is the Cultivation kind's PowerState: a character's standing at
// one power node, anchored to a Realm and the current Level within it, plus the
// points accumulated toward the next breakthrough. Power and lifespan are derived
// from the realm's ax+b formulas evaluated at Progress.
type CultivationState struct {
	Realm      Realm
	Level      Level   // current level within the realm
	Points     int     // accumulated toward Level.BreakthroughPoints (the breakthrough gate)
	Bottleneck int     // accumulated toward Level.BottleneckPoints (fills after breakthrough is full)
	Progress   float64 // x fed into the realm's power/lifespan formulas
}

// NewCultivationState validates and builds a CultivationState. It returns
// ErrInvalidProgress if progress is negative.
func NewCultivationState(realm Realm, level Level, points, bottleneck int, progress float64) (CultivationState, error) {
	if progress < 0 {
		return CultivationState{}, fmt.Errorf("%w: %g", ErrInvalidProgress, progress)
	}
	return CultivationState{Realm: realm, Level: level, Points: points, Bottleneck: bottleneck, Progress: progress}, nil
}

// Ready reports whether both gates of the current level are full — the character
// has broken through this level and is ready to advance.
func (c CultivationState) Ready() bool {
	return c.Points >= c.Level.BreakthroughPoints && c.Bottleneck >= c.Level.BottleneckPoints
}

// AdvanceWithin accumulates points into the state, filling the current level's
// breakthrough gate first and then, once it is full, its bottleneck gate. When
// both gates fill it breaks through to the next level in this realm (resetting
// the gates) and continues. It returns the updated state and any points left
// over after the realm's final level is fully broken through, for the caller to
// carry into the next realm.
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
		case c.Bottleneck < c.Level.BottleneckPoints:
			if n == 0 {
				return c, 0
			}
			add := min(n, c.Level.BottleneckPoints-c.Bottleneck)
			c.Bottleneck += add
			n -= add
		default: // both gates full — break through to the next level in this realm
			next, ok := c.Realm.nextLevel(c.Level.Number)
			if !ok {
				return c, n // realm ceiling reached; hand leftover to the caller
			}
			c.Level = next
			c.Points = 0
			c.Bottleneck = 0
		}
	}
}

// Kind implements PowerState.
func (CultivationState) Kind() PowerSystemType { return Cultivation }

// Power returns the cultivation's current power via the realm's formula.
func (c CultivationState) Power() float64 { return c.Realm.Power(c.Progress) }

// Lifespan returns the cultivation's current lifespan via the realm's formula.
func (c CultivationState) Lifespan() float64 { return c.Realm.Lifespan(c.Progress) }

var _ PowerState = CultivationState{}
