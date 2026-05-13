package rpg

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	// ErrInvalidQuestName is returned when a quest name is blank.
	ErrInvalidQuestName = errors.New("quest: name must not be empty")
	// ErrInvalidObjectiveDescription is returned for a blank objective.
	ErrInvalidObjectiveDescription = errors.New("quest: objective description must not be empty")
	// ErrObjectiveOrderExists is returned when an objective order is reused.
	ErrObjectiveOrderExists = errors.New("quest: objective order already used")
)

// Objective is a single ordered step of a Quest.
type Objective struct {
	Order       int
	Description string
}

// Quest is a unit of content with ordered objectives.
type Quest struct {
	Name        string
	Description string
	Objectives  []Objective
}

// NewQuest validates and builds a quest with no objectives.
func NewQuest(name, description string) (Quest, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Quest{}, ErrInvalidQuestName
	}
	return Quest{Name: name, Description: strings.TrimSpace(description)}, nil
}

// AddObjective inserts an ordered objective, rejecting a blank description or a
// duplicate order. Objectives are kept sorted by ascending Order.
func (q *Quest) AddObjective(order int, description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return ErrInvalidObjectiveDescription
	}
	for _, o := range q.Objectives {
		if o.Order == order {
			return fmt.Errorf("%w: %d", ErrObjectiveOrderExists, order)
		}
	}
	q.Objectives = append(q.Objectives, Objective{Order: order, Description: description})
	sort.Slice(q.Objectives, func(i, j int) bool { return q.Objectives[i].Order < q.Objectives[j].Order })
	return nil
}
