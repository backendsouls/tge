package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"tge/internal/adapter/sqlite"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

func TestAbilityRepository_SaveFindList(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "rpg.db")
	repo, err := sqlite.NewAbilityRepository(dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()

	a, _ := rpg.NewAbility("Berserk", "rage", "Common")
	if err := repo.Save(ctx, a); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := repo.Save(ctx, a); !errors.Is(err, port.ErrAbilityExists) {
		t.Fatalf("duplicate: got %v, want ErrAbilityExists", err)
	}
	got, err := repo.FindByName(ctx, "Berserk")
	if err != nil || got.Description != "rage" {
		t.Fatalf("find: %+v err=%v", got, err)
	}
	if _, err := repo.FindByName(ctx, "Nope"); !errors.Is(err, port.ErrAbilityNotFound) {
		t.Fatalf("missing: got %v", err)
	}
	list, err := repo.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %+v err=%v", list, err)
	}
}

func TestEquipmentRepository_RoundTripStats(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "rpg.db")
	repo, err := sqlite.NewEquipmentRepository(dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()

	eq, _ := rpg.NewEquipment("Steel Plate", rpg.Armor, rpg.Stats{STR: 2, VIT: 5})
	if err := repo.Save(ctx, eq); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := repo.FindByName(ctx, "Steel Plate")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Slot != rpg.Armor || got.Bonus.STR != 2 || got.Bonus.VIT != 5 {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestQuestRepository_Objectives(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "rpg.db")
	repo, err := sqlite.NewQuestRepository(dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()

	q, _ := rpg.NewQuest("Slay Dragon", "")
	if err := repo.Save(ctx, q); err != nil {
		t.Fatalf("save: %v", err)
	}
	// Insert out of order; the load must come back sorted by order.
	_ = repo.AddObjective(ctx, "Slay Dragon", rpg.Objective{Order: 2, Description: "Reach lair"})
	_ = repo.AddObjective(ctx, "Slay Dragon", rpg.Objective{Order: 1, Description: "Find map"})
	if err := repo.AddObjective(ctx, "Slay Dragon", rpg.Objective{Order: 1, Description: "dup"}); !errors.Is(err, rpg.ErrObjectiveOrderExists) {
		t.Fatalf("dup order: got %v", err)
	}
	got, err := repo.FindByName(ctx, "Slay Dragon")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got.Objectives) != 2 || got.Objectives[0].Order != 1 || got.Objectives[1].Order != 2 {
		t.Fatalf("objectives not sorted: %+v", got.Objectives)
	}
}

func TestCharacterRepository_StatsAndInventory(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "rpg.db")
	repo := newCharRepo(t, dsn)
	ctx := context.Background()

	c := mortalFixture(t, "Lin Feng", "MainCharacter", "Male", "Mortal Path")
	c.Class = rpg.Class{Name: "Warrior"}
	c.Stats = rpg.BaseStats()
	if err := repo.Save(ctx, c); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := repo.AddItem(ctx, "Lin Feng", "Gold", 10); err != nil {
		t.Fatalf("add item: %v", err)
	}
	if err := repo.AddItem(ctx, "Lin Feng", "Gold", 5); err != nil {
		t.Fatalf("add item: %v", err)
	}
	got, err := repo.FindByName(ctx, "Lin Feng")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Class.Name != "Warrior" || got.Stats.STR != 5 {
		t.Fatalf("class/stats not persisted: %+v", got)
	}
	if len(got.Inventory.Items) != 1 || got.Inventory.Items[0].Quantity != 15 {
		t.Fatalf("inventory not merged: %+v", got.Inventory.Items)
	}
}
