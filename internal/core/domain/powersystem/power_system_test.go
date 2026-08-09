package powersystem_test

import (
	"testing"

	"tge/internal/core/domain/powersystem"
)

func TestPowerSystem_Initialization(t *testing.T) {
	sys, err := powersystem.NewPowerSystem("Nasuverse", powersystem.Magic)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sys.Name != "Nasuverse" {
		t.Errorf("expected Nasuverse, got %s", sys.Name)
	}

	if sys.Nodes == nil {
		t.Errorf("expected Nodes map to be initialized")
	}
}

func TestPowerSystem_AddNode(t *testing.T) {
	sys, _ := powersystem.NewPowerSystem("Test", powersystem.Magic)
	node, _ := powersystem.NewPowerNode("Fireball", "Spell", nil)

	err := sys.AddNode(&node)
	if err != nil {
		t.Fatalf("expected no error adding node, got %v", err)
	}

	// Cannot add a node with the same ID twice
	err = sys.AddNode(&node)
	if err == nil {
		t.Errorf("expected error adding duplicate node")
	}
}

func TestPowerSystem_DAGValidation(t *testing.T) {
	sys, _ := powersystem.NewPowerSystem("MySystem", powersystem.Magic)
	nodeA, _ := powersystem.NewPowerNode("Node A", "Base", nil)
	nodeB, _ := powersystem.NewPowerNode("Node B", "Base", nil)

	_ = sys.AddNode(&nodeA)
	_ = sys.AddNode(&nodeB)

	// Add A as parent to B
	err := sys.AddEdge(nodeB.ID, nodeA.ID, powersystem.EdgeParent)
	if err != nil {
		t.Fatalf("expected no error adding valid edge: %v", err)
	}

	// Try to add B as parent to A, causing a cycle (A -> B -> A)
	err = sys.AddEdge(nodeA.ID, nodeB.ID, powersystem.EdgeParent)
	if err == nil {
		t.Errorf("expected error when adding cyclic edge, got nil")
	}
}
