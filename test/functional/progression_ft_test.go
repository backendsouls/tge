package functional_test

import (
	"context"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/character"
	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

// mock dependencies for the FT
type mockWorld struct{}

func (mockWorld) EnsureDefaults(ctx context.Context) (port.DefaultWorld, error) {
	return port.DefaultWorld{Species: character.Species{Name: "Human"}, PowerSystem: "HunterXHunter"}, nil
}

type mockSpecies struct{}

func (mockSpecies) FindByName(ctx context.Context, name string) (character.Species, error) {
	return character.Species{Name: name, Power: 1.0}, nil
}
func (mockSpecies) List(ctx context.Context) ([]character.Species, error) { return nil, nil }
func (mockSpecies) Save(ctx context.Context, s character.Species) error { return nil }

func TestFunctional_ProgressionDAG(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// 1. Boot up the repositories
	charRepo := file.NewCharacterRepository(tempDir + "/characters")
	sysRepo := file.NewPowerSystemRepository(tempDir + "/systems")

	// 2. Boot up the services
	sysSvc := service.NewPowerSystemService(sysRepo)
	charSvc := service.NewCharacterService(
		charRepo, sysRepo, mockSpecies{}, nil, nil, nil, mockWorld{}, service.CharacterDefaults{},
	)

	// 3. Define the Hunter x Hunter Nen System
	_, err := sysSvc.CreateSystem(ctx, "Nen", powersystem.Magic)
	if err != nil {
		t.Fatalf("failed to create system: %v", err)
	}

	// Add Nodes
	sysSvc.AddNode(ctx, port.AddNodeInput{System: "Nen", NodeID: "ten", Name: "Ten", Category: "Basic"})
	sysSvc.AddNode(ctx, port.AddNodeInput{System: "Nen", NodeID: "ren", Name: "Ren", Category: "Basic"})
	sysSvc.AddNode(ctx, port.AddNodeInput{System: "Nen", NodeID: "enhancement", Name: "Enhancement", Category: "Hatsu"})
	sysSvc.AddNode(ctx, port.AddNodeInput{System: "Nen", NodeID: "transmutation", Name: "Transmutation", Category: "Hatsu"})

	// Add Edges (Ren requires Ten)
	sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Nen", NodeID: "ren", TargetID: "ten", EdgeType: powersystem.EdgeParent})
	// Hatsu requires Ren
	sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Nen", NodeID: "enhancement", TargetID: "ren", EdgeType: powersystem.EdgeParent})
	sysSvc.AddEdge(ctx, port.AddEdgeInput{System: "Nen", NodeID: "transmutation", TargetID: "ren", EdgeType: powersystem.EdgeParent})

	// Make Enhancement and Transmutation mutually exclusive for this test scenario
	sys, _ := sysRepo.FindByName(ctx, "Nen")
	enhancementNode := sys.Nodes["enhancement"]
	enhancementNode.AddMutuallyExclusive("transmutation")
	transmutationNode := sys.Nodes["transmutation"]
	transmutationNode.AddMutuallyExclusive("enhancement")
	sysRepo.Save(ctx, sys)

	// 4. Create the Character
	charSvc.CreateCharacter(ctx, port.CreateCharacterInput{
		Name:    "Gon",
		Type:    string(character.MainCharacter),
		Gender:  string(character.Male),
		Species: "Human",
	})

	// 5. Attempt Invalid Progression (Missing Parent)
	// Try to learn Ren before Ten
	_, err = charSvc.TrainNode(ctx, port.TrainNodeInput{
		Character: "Gon",
		System:    "Nen",
		NodeID:    "ren",
	})
	if err == nil {
		t.Errorf("expected error when learning 'ren' without unlocking parent 'ten'")
	}

	// 6. Valid Progression (Learn Ten, then Ren, then Enhancement)
	_, err = charSvc.TrainNode(ctx, port.TrainNodeInput{Character: "Gon", System: "Nen", NodeID: "ten"})
	if err != nil {
		t.Errorf("failed to learn ten: %v", err)
	}
	_, err = charSvc.TrainNode(ctx, port.TrainNodeInput{Character: "Gon", System: "Nen", NodeID: "ren"})
	if err != nil {
		t.Errorf("failed to learn ren: %v", err)
	}
	gon, err := charSvc.TrainNode(ctx, port.TrainNodeInput{Character: "Gon", System: "Nen", NodeID: "enhancement"})
	if err != nil {
		t.Errorf("failed to learn enhancement: %v", err)
	}

	// 7. Verify Mutually Exclusive Blocking
	// Try to learn Transmutation while holding Enhancement
	_, err = charSvc.TrainNode(ctx, port.TrainNodeInput{
		Character: "Gon",
		System:    "Nen",
		NodeID:    "transmutation",
	})
	if err == nil {
		t.Errorf("expected error when trying to learn mutually exclusive node 'transmutation'")
	}

	// 8. Verify Power Scaling
	if len(gon.UnlockedNodes) != 3 {
		t.Errorf("expected 3 unlocked nodes, got %d", len(gon.UnlockedNodes))
	}
}
