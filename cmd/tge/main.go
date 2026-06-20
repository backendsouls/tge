// Command tge is the cultivation game CLI.
//
// main is the composition root of the hexagon: it is the only place that knows
// about concrete adapters. It wires the SQLite repository into the application
// service and hands that to the CLI adapter.
package main

import (
	"context"
	"fmt"
	"os"

	"tge/internal/adapter/cli"
	"tge/internal/adapter/sqlite"
	"tge/internal/config"
	"tge/internal/core/domain/character"
	"tge/internal/core/service"
)

func main() {
	os.Exit(run())
}

func run() int {
	dbPath := os.Getenv("TGE_DB")
	if dbPath == "" {
		dbPath = "tge.db"
	}

	defaults, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	realmRepo, err := sqlite.NewRealmRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer realmRepo.Close()

	psRepo, err := sqlite.NewPowerSystemRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer psRepo.Close()

	universeRepo, err := sqlite.NewUniverseRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer universeRepo.Close()

	charRepo, err := sqlite.NewCharacterRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer charRepo.Close()

	novelRepo, err := sqlite.NewNovelRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer novelRepo.Close()

	speciesRepo, err := sqlite.NewSpeciesRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer speciesRepo.Close()

	multiverseRepo, err := sqlite.NewMultiverseRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer multiverseRepo.Close()

	omniverseRepo, err := sqlite.NewOmniverseRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer omniverseRepo.Close()

	realityRepo, err := sqlite.NewRealityRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer realityRepo.Close()

	timelineRepo, err := sqlite.NewTimelineRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer timelineRepo.Close()

	abilityRepo, err := sqlite.NewAbilityRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer abilityRepo.Close()

	skillRepo, err := sqlite.NewSkillRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer skillRepo.Close()

	itemRepo, err := sqlite.NewItemRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer itemRepo.Close()

	effectRepo, err := sqlite.NewEffectRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer effectRepo.Close()

	equipmentRepo, err := sqlite.NewEquipmentRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer equipmentRepo.Close()

	professionRepo, err := sqlite.NewProfessionRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer professionRepo.Close()

	classRepo, err := sqlite.NewClassRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer classRepo.Close()

	questRepo, err := sqlite.NewQuestRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer questRepo.Close()

	recipeRepo, err := sqlite.NewRecipeRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer recipeRepo.Close()

	realmSvc := service.NewRealmService(realmRepo)
	psSvc := service.NewPowerSystemService(psRepo)
	universeSvc := service.NewUniverseService(universeRepo, psRepo, timelineRepo)
	multiverseSvc := service.NewMultiverseService(multiverseRepo, universeRepo, timelineRepo)
	omniverseSvc := service.NewOmniverseService(omniverseRepo, multiverseRepo, timelineRepo)
	realitySvc := service.NewRealityService(realityRepo, omniverseRepo, timelineRepo)
	timelineSvc := service.NewTimelineService(timelineRepo)
	speciesSvc := service.NewSpeciesService(speciesRepo)

	worldNames := service.WorldNames{
		Reality:     defaults.World.Reality,
		Omniverse:   defaults.World.Omniverse,
		Multiverse:  defaults.World.Multiverse,
		Universe:    defaults.World.Universe,
		Realm:       defaults.World.Realm,
		PowerSystem: defaults.World.PowerSystem,
	}
	humanBase := character.Species{
		Name:          defaults.Character.Species.Name,
		Power:         defaults.Character.Species.Power,
		Lifespan:      defaults.Character.Species.Lifespan,
		DefaultGender: character.Gender(defaults.Character.Species.DefaultGender),
	}
	charDefaults := service.CharacterDefaults{
		Gender: defaults.Character.Gender,
		Age:    defaults.Character.Age,
		Stats:  statsFrom(defaults.Character.Stats),
	}

	worldSvc := service.NewDefaultWorldService(realityRepo, omniverseRepo, multiverseRepo, universeRepo, psRepo, speciesRepo, timelineRepo, worldNames, humanBase)
	charSvc := service.NewCharacterService(charRepo, psRepo, speciesRepo, classRepo, professionRepo, itemRepo, worldSvc, charDefaults)
	novelSvc := service.NewNovelService(novelRepo, charRepo)

	rpgSvcs := cli.RPGServices{
		Abilities:   service.NewAbilityService(abilityRepo),
		Skills:      service.NewSkillService(skillRepo),
		Items:       service.NewItemService(itemRepo),
		Effects:     service.NewEffectService(effectRepo),
		Equipment:   service.NewEquipmentService(equipmentRepo),
		Professions: service.NewProfessionService(professionRepo),
		Classes:     service.NewClassService(classRepo),
		Quests:      service.NewQuestService(questRepo),
		Recipes:     service.NewRecipeService(recipeRepo),
	}

	// Seed the starter catalog (idempotent) before serving commands.
	if err := seedCatalog(context.Background(), seedDeps{
		species: speciesSvc,
		realms:  realmSvc,
		rpg:     rpgSvcs,
	}, defaults.Catalog); err != nil {
		fmt.Fprintf(os.Stderr, "error: seeding defaults: %v\n", err)
		return 1
	}

	app := cli.New(realmSvc, psSvc, universeSvc, multiverseSvc, omniverseSvc, realitySvc, timelineSvc, charSvc, novelSvc, speciesSvc, rpgSvcs, os.Stdout, os.Stderr)
	return app.Run(context.Background(), os.Args[1:])
}
