package powersystem_test

import (
	"testing"

	"tge/internal/core/domain/powersystem"
)

func TestPowerNode_Initialization(t *testing.T) {
	node, err := powersystem.NewPowerNode("Fireball", "Magic", []string{"spell", "fire"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Name != "Fireball" {
		t.Errorf("expected Name to be Fireball, got %s", node.Name)
	}

	if node.Category != "Magic" {
		t.Errorf("expected Category to be Magic, got %s", node.Category)
	}

	if len(node.Tags) != 2 || node.Tags[0] != "spell" {
		t.Errorf("expected tags [spell, fire], got %v", node.Tags)
	}

	// Test default stat vector initialization
	if node.StatVector == nil {
		t.Errorf("expected StatVector to be initialized to an empty map")
	}
}

func TestPowerNode_Relationships(t *testing.T) {
	node, _ := powersystem.NewPowerNode("Super Saiyan", "Transformation", nil)
	err := node.AddMutuallyExclusive("Kaioken")
	if err != nil {
		t.Fatalf("failed to add mutually exclusive node: %v", err)
	}

	if len(node.MutuallyExclusive) != 1 || node.MutuallyExclusive[0] != "Kaioken" {
		t.Errorf("expected Kaioken to be mutually exclusive")
	}

	// A node cannot be mutually exclusive with itself
	err = node.AddMutuallyExclusive(node.Name)
	if err == nil {
		t.Errorf("expected error when adding self as mutually exclusive, got nil")
	}
}
