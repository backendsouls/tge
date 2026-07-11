package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"tge/internal/adapter/cli"
	"tge/internal/config"
	"tge/internal/core/domain/cultivation"
	"tge/internal/core/domain/powersystem"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

// seedDeps are the driving ports the catalog seeder needs.
type seedDeps struct {
	species port.SpeciesService
	realms  port.RealmService
	systems port.PowerSystemService
	rpg     cli.RPGServices
}

// statsFrom converts a config stat block into the domain type.
func statsFrom(s config.Stats) rpg.Stats {
	return rpg.Stats{STR: s.STR, AGI: s.AGI, INT: s.INT, VIT: s.VIT, DEX: s.DEX, WIS: s.WIS, CHA: s.CHA, LUK: s.LUK}
}

// seedExists reports whether err is an "already exists" / duplicate outcome,
// which the idempotent seeder treats as success.
func seedExists(err error) bool {
	switch {
	case errors.Is(err, port.ErrSpeciesExists),
		errors.Is(err, port.ErrRealmExists),
		errors.Is(err, port.ErrClassExists),
		errors.Is(err, port.ErrProfessionExists),
		errors.Is(err, port.ErrItemExists),
		errors.Is(err, port.ErrAbilityExists),
		errors.Is(err, port.ErrSkillExists),
		errors.Is(err, port.ErrEffectExists),
		errors.Is(err, port.ErrEquipmentExists),
		errors.Is(err, port.ErrQuestExists),
		errors.Is(err, port.ErrRecipeExists),
		errors.Is(err, port.ErrPowerSystemExists),
		errors.Is(err, cultivation.ErrLevelNumberExists),
		errors.Is(err, powersystem.ErrNodeExists),
		errors.Is(err, rpg.ErrObjectiveOrderExists),
		errors.Is(err, rpg.ErrIngredientExists):
		return true
	}
	return false
}

// seed runs one create step, swallowing an "already exists" outcome so the
// catalog can be seeded repeatedly.
func seed(what string, err error) error {
	if err != nil && !seedExists(err) {
		return fmt.Errorf("seed %s: %w", what, err)
	}
	return nil
}

// seedCatalog idempotently creates every entity in the catalog. Existing
// entities are left untouched.
func seedCatalog(ctx context.Context, d seedDeps, cat config.Catalog) error {
	for _, sp := range cat.Species {
		if err := seed("species "+sp.Name, second(d.species.CreateSpecies(ctx, port.CreateSpeciesInput{
			Name: sp.Name, Power: sp.Power, Lifespan: sp.Lifespan, DefaultGender: sp.DefaultGender,
		}))); err != nil {
			return err
		}
	}

	for i, r := range cat.Realms {
		if err := seed("realm "+r.Name, second(d.realms.AddRealm(ctx, cultivation.RealmConfig{
			Name:                   r.Name,
			Tier:                   i + 1, // catalog order defines the realm sequence
			PowerMultiplier:        r.PowerMultiplier,
			PowerAdder:             r.PowerAdder,
			LifespanMultiplier:     r.LifespanMultiplier,
			LifespanAdder:          r.LifespanAdder,
			BottleneckPoints:       r.BottleneckPoints,
			MaxLevels:              r.MaxLevels,
			MainCharacterMaxLevels: r.MainCharacterMaxLevels,
		}))); err != nil {
			return err
		}
		for _, l := range realmLevels(r) {
			if err := seed(fmt.Sprintf("realm %s level %d", r.Name, l.Number),
				second(d.realms.AddLevel(ctx, port.AddLevelInput{
					Realm:              r.Name,
					Number:             l.Number,
					Name:               l.Name,
					BreakthroughPoints: l.BreakthroughPoints,
					BottleneckPoints:   l.BottleneckPoints,
				}))); err != nil {
				return err
			}
		}
	}

	for _, c := range cat.Classes {
		if err := seed("class "+c.Name, second(d.rpg.Classes.CreateClass(ctx, port.CreateClassInput{Name: c.Name, Description: c.Description}))); err != nil {
			return err
		}
	}
	for _, p := range cat.Professions {
		if err := seed("profession "+p.Name, second(d.rpg.Professions.CreateProfession(ctx, port.CreateProfessionInput{Name: p.Name, Description: p.Description}))); err != nil {
			return err
		}
	}
	for _, it := range cat.Items {
		if err := seed("item "+it.Name, second(d.rpg.Items.CreateItem(ctx, port.CreateItemInput{Name: it.Name, Description: it.Description}))); err != nil {
			return err
		}
	}
	for _, ab := range cat.Abilities {
		if err := seed("ability "+ab.Name, second(d.rpg.Abilities.CreateAbility(ctx, port.CreateAbilityInput{Name: ab.Name, Description: ab.Description}))); err != nil {
			return err
		}
	}
	for _, sk := range cat.Skills {
		if err := seed("skill "+sk.Name, second(d.rpg.Skills.CreateSkill(ctx, port.CreateSkillInput{Name: sk.Name, Description: sk.Description}))); err != nil {
			return err
		}
	}
	for _, e := range cat.Effects {
		if err := seed("effect "+e.Name, second(d.rpg.Effects.CreateEffect(ctx, port.CreateEffectInput{Name: e.Name, Kind: e.Kind, Description: e.Description}))); err != nil {
			return err
		}
	}
	for _, eq := range cat.Equipment {
		if err := seed("equipment "+eq.Name, second(d.rpg.Equipment.CreateEquipment(ctx, port.CreateEquipmentInput{Name: eq.Name, Slot: eq.Slot, Bonus: statsFrom(eq.Bonus)}))); err != nil {
			return err
		}
	}
	for _, q := range cat.Quests {
		if err := seed("quest "+q.Name, second(d.rpg.Quests.CreateQuest(ctx, port.CreateQuestInput{Name: q.Name, Description: q.Description}))); err != nil {
			return err
		}
		for _, o := range q.Objectives {
			if err := seed(fmt.Sprintf("quest %s objective %d", q.Name, o.Order),
				second(d.rpg.Quests.AddObjective(ctx, port.AddQuestObjectiveInput{Quest: q.Name, Order: o.Order, Description: o.Description}))); err != nil {
				return err
			}
		}
	}
	for _, rc := range cat.Recipes {
		if err := seed("recipe "+rc.Name, second(d.rpg.Recipes.CreateRecipe(ctx, port.CreateRecipeInput{Name: rc.Name, Output: rc.Output}))); err != nil {
			return err
		}
		for _, in := range rc.Inputs {
			if err := seed(fmt.Sprintf("recipe %s input %s", rc.Name, in.Item),
				second(d.rpg.Recipes.AddIngredient(ctx, port.AddIngredientInput{Recipe: rc.Name, Item: in.Item, Quantity: in.Quantity}))); err != nil {
				return err
			}
		}
	}

	for _, ps := range cat.PowerSystems {
		if err := seed("powersystem "+ps.Name, second(d.systems.CreateSystem(ctx, ps.Name, powersystem.PowerSystemType(ps.Kind)))); err != nil {
			return err
		}
		var seedPowers func(string, []config.PowerNode) error
		seedPowers = func(parentID string, powers []config.PowerNode) error {
			for _, p := range powers {
				nodeID := strings.ToLower(strings.ReplaceAll(p.Name, " ", "_"))
				// Add Node
				if err := seed("powersystem "+ps.Name+" node "+p.Name, second(d.systems.AddNode(ctx, port.AddNodeInput{
					System:   ps.Name,
					NodeID:   nodeID,
					Name:     p.Name,
					Category: "General",
				}))); err != nil {
					return err
				}

				// Add Edge if parent exists
				if parentID != "" {
					if err := seed("powersystem "+ps.Name+" edge "+p.Name, second(d.systems.AddEdge(ctx, port.AddEdgeInput{
						System:   ps.Name,
						NodeID:   nodeID,
						TargetID: parentID,
						EdgeType: powersystem.EdgeParent,
					}))); err != nil {
						return err
					}
				}

				if err := seedPowers(nodeID, p.Children); err != nil {
					return err
				}
			}
			return nil
		}
		if err := seedPowers("", ps.Powers); err != nil {
			return err
		}
	}

	return nil
}

// second discards the first (value) result of a create call, keeping the error,
// so create calls compose with seed().
func second[T any](_ T, err error) error { return err }

var ordinals = []string{
	"First", "Second", "Third", "Fourth", "Fifth", "Sixth", "Seventh", "Eighth", "Ninth", "Tenth",
}

// levelName returns an ordinal level name, e.g. 1 -> "First Level".
func levelName(n int) string {
	if n >= 1 && n <= len(ordinals) {
		return ordinals[n-1] + " Level"
	}
	return fmt.Sprintf("Level %d", n)
}

// realmLevels returns the realm's explicit levels, or, when none are given,
// LevelCount ordinally-named levels.
func realmLevels(r config.Realm) []config.Level {
	if len(r.Levels) > 0 {
		return r.Levels
	}
	levels := make([]config.Level, 0, r.LevelCount)
	for i := 1; i <= r.LevelCount; i++ {
		levels = append(levels, config.Level{
			Number:             i,
			Name:               levelName(i),
			BreakthroughPoints: r.LevelBreakthroughPoints,
			BottleneckPoints:   r.LevelBottleneckPoints,
		})
	}
	return levels
}
