package functional_test

import (
	"context"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/character"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

type mockItem struct{}

func (mockItem) FindByName(ctx context.Context, name string) (rpg.Item, error) {
	return rpg.Item{Name: name}, nil
}
func (mockItem) List(ctx context.Context) ([]rpg.Item, error) { return nil, nil }
func (mockItem) Save(ctx context.Context, i rpg.Item) error { return nil }

func TestFunctional_CharacterInventoryAndState(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	charRepo := file.NewCharacterRepository(tempDir + "/characters")

	charSvc := service.NewCharacterService(
		charRepo, nil, mockSpecies{}, nil, nil, mockItem{}, mockWorld{}, service.CharacterDefaults{},
	)

	// Create Character
	_, err := charSvc.CreateCharacter(ctx, port.CreateCharacterInput{
		Name:    "Kirito",
		Type:    string(character.MainCharacter),
		Gender:  string(character.Male),
		Species: "Human",
	})
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	// Give Items
	charSvc.GiveItem(ctx, port.GiveItemInput{Character: "Kirito", Item: "Elucidator", Quantity: 1})
	charSvc.GiveItem(ctx, port.GiveItemInput{Character: "Kirito", Item: "Healing Potion", Quantity: 5})

	// Stack existing items
	charSvc.GiveItem(ctx, port.GiveItemInput{Character: "Kirito", Item: "Healing Potion", Quantity: 5})

	// Verify State Persistence
	c, err := charRepo.FindByName(ctx, "Kirito")
	if err != nil {
		t.Fatalf("failed to load character: %v", err)
	}

	if len(c.Inventory.Items) != 2 {
		t.Errorf("expected 2 distinct items, got %d", len(c.Inventory.Items))
	}

	for _, item := range c.Inventory.Items {
		if item.Item == "Healing Potion" && item.Quantity != 10 {
			t.Errorf("expected 10 Healing Potions, got %d", item.Quantity)
		}
		if item.Item == "Elucidator" && item.Quantity != 1 {
			t.Errorf("expected 1 Elucidator, got %d", item.Quantity)
		}
	}

	// Verify MechanicState defaults
	if c.MechanicState.Tier != 0 {
		t.Errorf("expected default tier 0, got %d", c.MechanicState.Tier)
	}
	if c.CalculateTotalPower() != 1.0 {
		t.Errorf("expected default power 1.0, got %f", c.CalculateTotalPower())
	}
}
