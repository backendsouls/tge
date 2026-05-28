package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cosmology"
)

var (
	ErrTimelineNotFound = errors.New("timeline: not found")
	ErrTimelineExists   = errors.New("timeline: already exists")
)

// LocationKind identifies which kind of location owns a timeline.
type LocationKind string

const (
	LocationBox        LocationKind = "box"
	LocationOmniverse  LocationKind = "omniverse"
	LocationMultiverse LocationKind = "multiverse"
	LocationUniverse   LocationKind = "universe"
	LocationRealm      LocationKind = "realm"
)

// Valid reports whether k is a known location kind.
func (k LocationKind) Valid() bool {
	switch k {
	case LocationBox, LocationOmniverse, LocationMultiverse, LocationUniverse, LocationRealm:
		return true
	}
	return false
}

// LocationRef identifies the location that owns a timeline. Universe scopes a
// realm, whose name is unique only within its universe; it is empty for the
// other kinds, whose names are globally unique.
type LocationRef struct {
	Kind     LocationKind
	Name     string
	Universe string // realm only
}

// AddTimelineEventInput describes an event to append to a location's timeline.
type AddTimelineEventInput struct {
	Owner       LocationRef
	Order       int
	Description string
}

// TimelineRepository persists each location's single timeline and its events.
type TimelineRepository interface {
	// Save persists a location's timeline (without events), 1:1 with the owner.
	Save(ctx context.Context, owner LocationRef, t cosmology.Timeline) error
	Find(ctx context.Context, owner LocationRef) (cosmology.Timeline, error)
	AddEvent(ctx context.Context, owner LocationRef, e cosmology.Event) error
}

// TimelineService is a driving port for viewing and extending a location's timeline.
type TimelineService interface {
	GetTimeline(ctx context.Context, owner LocationRef) (cosmology.Timeline, error)
	AddEvent(ctx context.Context, in AddTimelineEventInput) (cosmology.Timeline, error)
}
