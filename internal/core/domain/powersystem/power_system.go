package powersystem

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSystemName = errors.New("power system: name must not be empty")
	ErrNodeExists        = errors.New("power system: node already exists")
	ErrNodeNotFound      = errors.New("power system: node not found")
	ErrCyclicDependency  = errors.New("power system: cyclic dependency detected")
)

type PowerSystemType string

const (
	Cultivation PowerSystemType = "Cultivation"
	Magic       PowerSystemType = "Magic"
	SuperPower  PowerSystemType = "SuperPower"
)

func (k PowerSystemType) Valid() bool {
	switch k {
	case Cultivation, Magic, SuperPower:
		return true
	}
	return false
}

type EdgeType string

const (
	EdgeParent            EdgeType = "parent"
	EdgeSibling           EdgeType = "sibling"
	EdgeMutuallyExclusive EdgeType = "mutually_exclusive"
)

// PowerSystem is a named DAG (Directed Acyclic Graph) of PowerNodes.
// It stores nodes in a flat map for O(1) traversal and relationship mapping.
type PowerSystem struct {
	Name            string
	PowerSystemType PowerSystemType
	Nodes           map[string]*PowerNode
}

// NewPowerSystem validates and builds an empty power system.
func NewPowerSystem(name string, kind PowerSystemType) (PowerSystem, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return PowerSystem{}, ErrInvalidSystemName
	}
	if string(kind) == "" {
		kind = Cultivation
	}
	if !kind.Valid() {
		return PowerSystem{}, fmt.Errorf("power system: invalid kind %q", kind)
	}
	return PowerSystem{
		Name:            name,
		PowerSystemType: kind,
		Nodes:           make(map[string]*PowerNode),
	}, nil
}

// AddNode adds a PowerNode to the system graph.
func (ps *PowerSystem) AddNode(node *PowerNode) error {
	if _, exists := ps.Nodes[node.ID]; exists {
		return fmt.Errorf("%w: %s", ErrNodeExists, node.ID)
	}
	ps.Nodes[node.ID] = node
	return nil
}

// AddEdge registers a relationship between two nodes in the system.
// If EdgeParent is used, it verifies that the relationship does not introduce a cycle.
func (ps *PowerSystem) AddEdge(nodeID, targetID string, edgeType EdgeType) error {
	node, ok := ps.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, nodeID)
	}
	_, ok = ps.Nodes[targetID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, targetID)
	}

	if edgeType == EdgeSibling {
		for _, s := range node.Siblings {
			if s == targetID {
				return nil
			}
		}
		node.Siblings = append(node.Siblings, targetID)
		return nil
	}

	if edgeType == EdgeParent {
		for _, p := range node.Parents {
			if p == targetID {
				return nil
			}
		}
		// Cycle Detection via DFS
		if ps.wouldCreateCycle(nodeID, targetID) {
			return fmt.Errorf("%w: %s -> %s", ErrCyclicDependency, nodeID, targetID)
		}
		node.Parents = append(node.Parents, targetID)
	}

	if edgeType == EdgeMutuallyExclusive {
		_ = node.AddMutuallyExclusive(targetID)
		targetNode, _ := ps.Nodes[targetID]
		_ = targetNode.AddMutuallyExclusive(nodeID)
	}

	return nil
}

// wouldCreateCycle runs a Depth First Search (DFS) to see if adding `targetID` as a parent of `nodeID`
// would allow a path from `targetID` back to `nodeID`.
func (ps *PowerSystem) wouldCreateCycle(nodeID, targetID string) bool {
	if nodeID == targetID {
		return true // A node cannot be its own parent
	}

	visited := make(map[string]bool)
	var dfs func(current string) bool
	dfs = func(current string) bool {
		if current == nodeID {
			return true // Found a path back to the starting node!
		}
		if visited[current] {
			return false
		}
		visited[current] = true

		currNode, ok := ps.Nodes[current]
		if !ok {
			return false
		}
		for _, parent := range currNode.Parents {
			if dfs(parent) {
				return true
			}
		}
		return false
	}

	// If we trace upward from the proposed parent, do we eventually hit nodeID?
	return dfs(targetID)
}

// Names returns every node ID in the system.
func (ps *PowerSystem) Names() []string {
	var out []string
	for id := range ps.Nodes {
		out = append(out, id)
	}
	return out
}
