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
	Realm    Realm
	Level    Level   // current level within the realm
	Points   int     // accumulated toward Level.BreakthroughPoints
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

// Kind implements PowerState.
func (CultivationState) Kind() SystemKind { return Cultivation }

// Power returns the cultivation's current power via the realm's formula.
func (c CultivationState) Power() float64 { return c.Realm.Power(c.Progress) }

// Lifespan returns the cultivation's current lifespan via the realm's formula.
func (c CultivationState) Lifespan() float64 { return c.Realm.Lifespan(c.Progress) }

var _ PowerState = CultivationState{}
