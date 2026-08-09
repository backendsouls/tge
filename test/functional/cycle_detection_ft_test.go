package functional_test

import (
	"context"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

func TestFunctional_CycleDetectionFT(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	sysRepo := file.NewPowerSystemRepository(tempDir + "/systems")
	sysSvc := service.NewPowerSystemService(sysRepo)

	_, _ = sysSvc.CreateSystem(ctx, "Alchemy", powersystem.Magic)
	_, _ = sysSvc.AddNode(ctx, port.AddNodeInput{System: "Alchemy", NodeID: "comprehension", Name: "Comprehension"})
	_, _ = sysSvc.AddNode(ctx, port.AddNodeInput{System: "Alchemy", NodeID: "deconstruction", Name: "Deconstruction"})
	_, _ = sysSvc.AddNode(ctx, port.AddNodeInput{System: "Alchemy", NodeID: "reconstruction", Name: "Reconstruction"})

	// Create valid chain: Comprehension -> Deconstruction -> Reconstruction
	_, err := sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Alchemy", NodeID: "deconstruction", TargetID: "comprehension", EdgeType: powersystem.EdgeParent})
	if err != nil {
		t.Fatalf("unexpected error on valid edge: %v", err)
	}

	_, err = sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Alchemy", NodeID: "reconstruction", TargetID: "deconstruction", EdgeType: powersystem.EdgeParent})
	if err != nil {
		t.Fatalf("unexpected error on valid edge: %v", err)
	}

	// Attempt to create an invalid cycle: Reconstruction -> Comprehension
	// This would mean Comprehension requires Reconstruction, which requires Deconstruction, which requires Comprehension.
	_, err = sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Alchemy", NodeID: "comprehension", TargetID: "reconstruction", EdgeType: powersystem.EdgeParent})

	if err == nil {
		t.Fatalf("expected cycle detection to prevent recursive edge, got nil")
	}

	// Verify the DAG was not corrupted
	sys, _ := sysRepo.FindByName(ctx, "Alchemy")
	compNode := sys.Nodes["comprehension"]
	if len(compNode.Parents) != 0 {
		t.Errorf("expected Comprehension to have 0 parents, got %d", len(compNode.Parents))
	}
}
