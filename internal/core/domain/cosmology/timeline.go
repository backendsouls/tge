package cosmology

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrInvalidTimelineName     = errors.New("timeline: name must not be empty")
	ErrInvalidEventDescription = errors.New("timeline: event description must not be empty")
	ErrEventOrderExists        = errors.New("timeline: event order already used")
)

// Event is a single ordered point on a Timeline.
type Event struct {
	Order       int
	Description string
}

// Timeline is the ordered sequence of events owned by a location (a realm,
// universe, multiverse, omniverse or box). Every location has exactly one.
type Timeline struct {
	Name   string
	Events []Event
}

// NewTimeline validates and creates a new, empty Timeline.
func NewTimeline(name string) (Timeline, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Timeline{}, ErrInvalidTimelineName
	}
	return Timeline{Name: name}, nil
}

// AddEvent inserts an ordered event, rejecting a blank description or a
// duplicate order. Events are kept sorted by ascending Order.
func (t *Timeline) AddEvent(order int, description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return ErrInvalidEventDescription
	}
	for _, e := range t.Events {
		if e.Order == order {
			return fmt.Errorf("%w: %d", ErrEventOrderExists, order)
		}
	}
	t.Events = append(t.Events, Event{Order: order, Description: description})
	sort.Slice(t.Events, func(i, j int) bool { return t.Events[i].Order < t.Events[j].Order })
	return nil
}
