package powersystem

import (
	"errors"
	"strings"
)

var (
	ErrInvalidPowerNodeName  = errors.New("power node: name must not be empty")
	ErrSelfMutuallyExclusive = errors.New("power node: cannot be mutually exclusive with itself")
)

// PowerNode replaces the old hierarchical Power tree. It operates as a node within
// a Directed Acyclic Graph (DAG), enabling multi-parent and mutually exclusive relationships.
type PowerNode struct {
	ID                string
	Name              string
	Category          string
	Tags              []string
	Parents           []string
	Siblings          []string
	MutuallyExclusive []string
	BasePower         float64
	StatVector        map[string]float64
	MaterialReq       map[string]int
	Drawbacks         []string
}

// NewPowerNode validates and creates a new PowerNode with default mappings.
func NewPowerNode(name, category string, tags []string) (PowerNode, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return PowerNode{}, ErrInvalidPowerNodeName
	}

	// Default empty slices/maps so they aren't nil
	if tags == nil {
		tags = []string{}
	}

	return PowerNode{
		ID:                strings.ToLower(strings.ReplaceAll(name, " ", "_")),
		Name:              name,
		Category:          category,
		Tags:              tags,
		Parents:           []string{},
		Siblings:          []string{},
		MutuallyExclusive: []string{},
		StatVector:        make(map[string]float64),
		MaterialReq:       make(map[string]int),
		Drawbacks:         []string{},
	}, nil
}

// AddMutuallyExclusive registers another node's ID as mutually exclusive.
// Returns an error if attempting to add itself.
func (p *PowerNode) AddMutuallyExclusive(otherID string) error {
	if p.Name == otherID || p.ID == otherID {
		return ErrSelfMutuallyExclusive
	}
	for _, id := range p.MutuallyExclusive {
		if id == otherID {
			return nil
		}
	}
	p.MutuallyExclusive = append(p.MutuallyExclusive, otherID)
	return nil
}
