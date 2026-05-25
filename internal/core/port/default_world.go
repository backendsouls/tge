package port

import (
	"context"
	"tge/internal/core/domain/character"
)

// DefaultWorld names the entities of the default cosmology that backs a
// name-only main character: the full Box -> Omniverse -> Multiverse -> Universe
// -> Realm chain, the default power system, and the base species used as the
// character's initial-status template.
type DefaultWorld struct {
	Reality     string
	Omniverse   string
	Multiverse  string
	Universe    string
	Realm       string
	PowerSystem string
	Species     character.Species
}

// DefaultWorldProvisioner ensures the default cosmology exists, creating it on
// first use and returning the names to attach a new main character to. It is
// idempotent: calling it repeatedly leaves the same single default world.
type DefaultWorldProvisioner interface {
	EnsureDefaults(ctx context.Context) (DefaultWorld, error)
}
