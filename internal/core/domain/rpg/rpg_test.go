package rpg_test

import (
	"testing"

	"tge/internal/core/domain/rpg"
)

func TestNewStats(t *testing.T) {
	if _, err := rpg.NewStats(1, 1, 1, 1, 1, 1, 1, -1); err != rpg.ErrNegativeStat {
		t.Errorf("expected ErrNegativeStat, got %v", err)
	}
	s, err := rpg.NewStats(1, 2, 3, 4, 5, 6, 7, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := s.Add(rpg.BaseStats())
	if sum.STR != 6 || sum.LUK != 13 {
		t.Errorf("Add mismatch: %+v", sum)
	}
}

func TestNamedEntities(t *testing.T) {
	if _, err := rpg.NewAbility("  ", ""); err != rpg.ErrInvalidAbilityName {
		t.Errorf("ability: got %v", err)
	}
	if _, err := rpg.NewSkill("", ""); err != rpg.ErrInvalidSkillName {
		t.Errorf("skill: got %v", err)
	}
	if _, err := rpg.NewItem("", ""); err != rpg.ErrInvalidItemName {
		t.Errorf("item: got %v", err)
	}
	if _, err := rpg.NewProfession("", ""); err != rpg.ErrInvalidProfessionName {
		t.Errorf("profession: got %v", err)
	}
	if _, err := rpg.NewClass("", ""); err != rpg.ErrInvalidClassName {
		t.Errorf("class: got %v", err)
	}
	a, err := rpg.NewAbility(" Fireball ", " burns ")
	if err != nil || a.Name != "Fireball" || a.Description != "burns" {
		t.Errorf("ability trim: %+v err=%v", a, err)
	}
}

func TestNewEffect(t *testing.T) {
	if _, err := rpg.NewEffect("Poison", "Nope", ""); err == nil {
		t.Error("expected invalid kind error")
	}
	e, err := rpg.NewEffect("Poison", rpg.Status, "damage over time")
	if err != nil || e.Kind != rpg.Status {
		t.Errorf("effect: %+v err=%v", e, err)
	}
}

func TestNewEquipment(t *testing.T) {
	if _, err := rpg.NewEquipment("Sword", "Hand", rpg.Stats{}); err == nil {
		t.Error("expected invalid slot error")
	}
	eq, err := rpg.NewEquipment("Sword", rpg.Weapon, rpg.Stats{STR: 3})
	if err != nil || eq.Slot != rpg.Weapon || eq.Bonus.STR != 3 {
		t.Errorf("equipment: %+v err=%v", eq, err)
	}
}

func TestQuest_AddObjective(t *testing.T) {
	q, _ := rpg.NewQuest("Slay the Dragon", "")
	if err := q.AddObjective(1, ""); err != rpg.ErrInvalidObjectiveDescription {
		t.Errorf("got %v", err)
	}
	_ = q.AddObjective(2, "Reach the lair")
	_ = q.AddObjective(1, "Find the map")
	if len(q.Objectives) != 2 || q.Objectives[0].Order != 1 {
		t.Fatalf("objectives not sorted: %+v", q.Objectives)
	}
	if err := q.AddObjective(1, "dup"); err == nil {
		t.Error("expected duplicate order error")
	}
}

func TestRecipe_AddInput(t *testing.T) {
	if _, err := rpg.NewRecipe("Potion", ""); err != rpg.ErrInvalidRecipeOutput {
		t.Errorf("got %v", err)
	}
	r, _ := rpg.NewRecipe("Potion", "Health Potion")
	if err := r.AddInput("Herb", 0); err != rpg.ErrInvalidIngredientQuantity {
		t.Errorf("got %v", err)
	}
	_ = r.AddInput("Herb", 2)
	if err := r.AddInput("Herb", 1); err == nil {
		t.Error("expected duplicate ingredient error")
	}
}

func TestInventory_Add(t *testing.T) {
	var inv rpg.Inventory
	if err := inv.Add("", 1); err != rpg.ErrInvalidInventoryItem {
		t.Errorf("got %v", err)
	}
	_ = inv.Add("Gold", 10)
	_ = inv.Add("Gold", 5)
	if len(inv.Items) != 1 || inv.Items[0].Quantity != 15 {
		t.Fatalf("stacks not merged: %+v", inv.Items)
	}
}
