package service

import (
	"context"
	"errors"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/cosmology"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

// Names of the entities making up the default cosmology, used as the built-in
// fallback when configuration does not supply them. A name-only main character
// is born into this world.
const (
	DefaultReality     = "The Box"
	DefaultOmniverse   = "Origin Omniverse"
	DefaultMultiverse  = "Origin Multiverse"
	DefaultUniverse    = "Origin Universe"
	DefaultRealm       = "Mortal Realm"
	DefaultPowerSystem = "Mortal Path"
)

// WorldNames names the entities of the default cosmology. Empty fields fall back
// to the built-in Default* names.
type WorldNames struct {
	Reality     string
	Omniverse   string
	Multiverse  string
	Universe    string
	Realm       string
	PowerSystem string
}

func (n WorldNames) withDefaults() WorldNames {
	if n.Reality == "" {
		n.Reality = DefaultReality
	}
	if n.Omniverse == "" {
		n.Omniverse = DefaultOmniverse
	}
	if n.Multiverse == "" {
		n.Multiverse = DefaultMultiverse
	}
	if n.Universe == "" {
		n.Universe = DefaultUniverse
	}
	if n.Realm == "" {
		n.Realm = DefaultRealm
	}
	if n.PowerSystem == "" {
		n.PowerSystem = DefaultPowerSystem
	}
	return n
}

// DefaultWorldService provisions the default cosmology (Box -> Omniverse ->
// Multiverse -> Universe -> Realm), a default power system, and the Human base
// species. It implements port.DefaultWorldProvisioner and is idempotent: each
// step ignores an "already exists" outcome, so repeated calls converge on a
// single default world.
type DefaultWorldService struct {
	realities   port.RealityRepository
	omniverses  port.OmniverseRepository
	multiverses port.MultiverseRepository
	universes   port.UniverseRepository
	systems     port.PowerSystemRepository
	species     port.SpeciesRepository
	timelines   port.TimelineRepository
	names       WorldNames
	human       character.Species
}

// NewDefaultWorldService wires the provisioner to the driven repositories it
// needs to build the default cosmology. names and human supply the configurable
// default cosmology names and Human base species; empty/zero values fall back to
// the built-in defaults.
func NewDefaultWorldService(
	realities port.RealityRepository,
	omniverses port.OmniverseRepository,
	multiverses port.MultiverseRepository,
	universes port.UniverseRepository,
	systems port.PowerSystemRepository,
	species port.SpeciesRepository,
	timelines port.TimelineRepository,
	names WorldNames,
	human character.Species,
) *DefaultWorldService {
	if human.Name == "" {
		human = character.HumanBase()
	}
	return &DefaultWorldService{
		realities:   realities,
		omniverses:  omniverses,
		multiverses: multiverses,
		universes:   universes,
		systems:     systems,
		species:     species,
		timelines:   timelines,
		names:       names.withDefaults(),
		human:       human,
	}
}

// EnsureDefaults creates any missing parts of the default world and returns the
// names a new main character should be attached to.
func (s *DefaultWorldService) EnsureDefaults(ctx context.Context) (port.DefaultWorld, error) {
	n := s.names
	if err := ignoreExists(s.species.Save(ctx, s.human), port.ErrSpeciesExists); err != nil {
		return port.DefaultWorld{}, err
	}

	if err := ignoreExists(s.systems.Create(ctx, n.PowerSystem), port.ErrPowerSystemExists); err != nil {
		return port.DefaultWorld{}, err
	}

	if err := ignoreExists(s.universes.Create(ctx, n.Universe), port.ErrUniverseExists); err != nil {
		return port.DefaultWorld{}, err
	}
	// SaveSystems/SaveRealms replace the universe's membership, so setting the
	// fixed defaults is itself idempotent.
	u := cosmology.Universe{
		Name:    n.Universe,
		Systems: []progression.PowerSystem{{Name: n.PowerSystem}},
		Realms:  []cosmology.Location{{Name: n.Realm}},
	}
	if err := s.universes.SaveSystems(ctx, u); err != nil {
		return port.DefaultWorld{}, err
	}
	if err := s.universes.SaveRealms(ctx, u); err != nil {
		return port.DefaultWorld{}, err
	}

	if err := ignoreExists(s.multiverses.Save(ctx, cosmology.Multiverse{Name: n.Multiverse}), port.ErrMultiverseExists); err != nil {
		return port.DefaultWorld{}, err
	}
	if err := ignoreExists(s.multiverses.AddUniverse(ctx, n.Multiverse, n.Universe), cosmology.ErrMultiverseUniverseExists); err != nil {
		return port.DefaultWorld{}, err
	}

	if err := ignoreExists(s.omniverses.Save(ctx, cosmology.Omniverse{Name: n.Omniverse}), port.ErrOmniverseExists); err != nil {
		return port.DefaultWorld{}, err
	}
	if err := ignoreExists(s.omniverses.AddMultiverse(ctx, n.Omniverse, n.Multiverse), cosmology.ErrOmniverseMultiverseExists); err != nil {
		return port.DefaultWorld{}, err
	}

	if err := ignoreExists(s.realities.Save(ctx, cosmology.Reality{Name: n.Reality}), port.ErrRealityExists); err != nil {
		return port.DefaultWorld{}, err
	}
	if err := ignoreExists(s.realities.AddOmniverse(ctx, n.Reality, n.Omniverse), cosmology.ErrRealityOmniverseExists); err != nil {
		return port.DefaultWorld{}, err
	}

	// Every location in the default cosmology owns a timeline.
	for _, owner := range []port.LocationRef{
		{Kind: port.LocationBox, Name: n.Reality},
		{Kind: port.LocationOmniverse, Name: n.Omniverse},
		{Kind: port.LocationMultiverse, Name: n.Multiverse},
		{Kind: port.LocationUniverse, Name: n.Universe},
		{Kind: port.LocationRealm, Name: n.Realm, Universe: n.Universe},
	} {
		if err := ensureTimeline(ctx, s.timelines, owner); err != nil {
			return port.DefaultWorld{}, err
		}
	}

	return port.DefaultWorld{
		Reality:     n.Reality,
		Omniverse:   n.Omniverse,
		Multiverse:  n.Multiverse,
		Universe:    n.Universe,
		Realm:       n.Realm,
		PowerSystem: n.PowerSystem,
		Species:     s.human,
	}, nil
}

// ignoreExists returns nil when err is the given "already exists" sentinel,
// making a create step idempotent, and returns err otherwise.
func ignoreExists(err, exists error) error {
	if errors.Is(err, exists) {
		return nil
	}
	return err
}

var _ port.DefaultWorldProvisioner = (*DefaultWorldService)(nil)
