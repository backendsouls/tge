package service_test

import (
	"context"
	"errors"
	"testing"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

type stubCharRepo struct {
	mainChars    []character.Character
	byName       map[string]character.Character
	saveErr      error
	saved        []character.Character
	cultivations []port.CultivationRecord
}

func (s *stubCharRepo) Save(_ context.Context, c character.Character) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved = append(s.saved, c)
	return nil
}

func (s *stubCharRepo) FindByName(_ context.Context, name string) (character.Character, error) {
	c, ok := s.byName[name]
	if !ok {
		return character.Character{}, port.ErrCharacterNotFound
	}
	return c, nil
}

func (s *stubCharRepo) MainCharacters(context.Context) ([]character.Character, error) {
	return s.mainChars, nil
}

func (s *stubCharRepo) List(context.Context) ([]character.Character, error) {
	return s.saved, nil
}

func (s *stubCharRepo) AddItem(context.Context, string, string, int) error { return nil }

func (s *stubCharRepo) SaveCultivation(_ context.Context, _ string, rec port.CultivationRecord) error {
	s.cultivations = append(s.cultivations, rec)
	return nil
}

// stubClassRepo, stubProfessionRepo and stubItemRepo are permissive RPG-repo
// doubles: any looked-up name resolves, so optional class/profession validation
// passes when those fields are set.
type stubClassRepo struct{}

func (stubClassRepo) Save(context.Context, rpg.Class) error { return nil }
func (stubClassRepo) FindByName(_ context.Context, name string) (rpg.Class, error) {
	return rpg.Class{Name: name}, nil
}
func (stubClassRepo) List(context.Context) ([]rpg.Class, error) { return nil, nil }

type stubProfessionRepo struct{}

func (stubProfessionRepo) Save(context.Context, rpg.Profession) error { return nil }
func (stubProfessionRepo) FindByName(_ context.Context, name string) (rpg.Profession, error) {
	return rpg.Profession{Name: name}, nil
}
func (stubProfessionRepo) List(context.Context) ([]rpg.Profession, error) { return nil, nil }

type stubItemRepo struct{}

func (stubItemRepo) Save(context.Context, rpg.Item) error { return nil }
func (stubItemRepo) FindByName(_ context.Context, name string) (rpg.Item, error) {
	return rpg.Item{Name: name}, nil
}
func (stubItemRepo) List(context.Context) ([]rpg.Item, error) { return nil, nil }

// stubSystemRepo resolves a known set of power systems by name.
type stubSystemRepo struct {
	systems map[string]progression.PowerSystem
}

func (s *stubSystemRepo) Create(context.Context, string) error                      { return nil }
func (s *stubSystemRepo) List(context.Context) ([]progression.PowerSystem, error)   { return nil, nil }
func (s *stubSystemRepo) SavePowers(context.Context, progression.PowerSystem) error { return nil }
func (s *stubSystemRepo) FindByName(_ context.Context, name string) (progression.PowerSystem, error) {
	ps, ok := s.systems[name]
	if !ok {
		return progression.PowerSystem{}, port.ErrPowerSystemNotFound
	}
	return ps, nil
}

func knownSystems(names ...string) *stubSystemRepo {
	m := map[string]progression.PowerSystem{}
	for _, n := range names {
		m[n] = progression.PowerSystem{Name: n}
	}
	return &stubSystemRepo{systems: m}
}

func TestCharacterService_CreateCharacter(t *testing.T) {
	systems := knownSystems("Universe A Cultivation", "Universe B Sorcery")

	t.Run("creates a main character and persists it", func(t *testing.T) {
		chars := &stubCharRepo{}
		svc := service.NewCharacterService(chars, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})

		got, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name:    "Lin Feng",
			Type:    "MainCharacter",
			Gender:  "Female",
			Systems: []string{"Universe A Cultivation", "Universe B Sorcery"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != character.MainCharacter || len(got.PowerSystems) != 2 {
			t.Errorf("got %+v", got)
		}
		if len(chars.saved) != 1 {
			t.Fatal("character not persisted")
		}
	})

	t.Run("name-only main character is born from the Human base into the default world", func(t *testing.T) {
		chars := &stubCharRepo{}
		world := &stubWorld{}
		svc := service.NewCharacterService(chars, knownSystems(service.DefaultPowerSystem), &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, world, service.CharacterDefaults{})

		got, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "Lin Feng", Type: "MainCharacter",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !world.called {
			t.Error("default world was not provisioned")
		}
		if got.Gender != character.Male {
			t.Errorf("Gender = %q, want default Male", got.Gender)
		}
		if len(got.Species) != 1 || got.Species[0].Name != character.HumanBaseName {
			t.Errorf("Species = %+v, want %q", got.Species, character.HumanBaseName)
		}
		if len(got.PowerSystems) != 1 || got.PowerSystems[0].Name != service.DefaultPowerSystem {
			t.Errorf("PowerSystems = %+v, want the default power system", got.PowerSystems)
		}
	})

	t.Run("allows multiple main characters", func(t *testing.T) {
		chars := &stubCharRepo{mainChars: []character.Character{
			{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female},
		}}
		svc := service.NewCharacterService(chars, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})

		if _, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "Mu Chen", Type: "MainCharacter", Gender: "Male", Systems: []string{"Universe A Cultivation"},
		}); err != nil {
			t.Fatalf("second main character rejected: %v", err)
		}
		if len(chars.saved) != 1 {
			t.Fatal("second main character not persisted")
		}
	})

	t.Run("Hero allowed only when a female main character exists", func(t *testing.T) {
		femaleMCs := []character.Character{{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female}}
		svc := service.NewCharacterService(&stubCharRepo{mainChars: femaleMCs}, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})
		if _, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "Yun Hai", Type: "Hero", Gender: "Male", Systems: []string{"Universe A Cultivation"},
		}); err != nil {
			t.Fatalf("Hero with female MC: %v", err)
		}

		maleMCs := []character.Character{{Name: "Mu Chen", Type: character.MainCharacter, Gender: character.Male}}
		svc = service.NewCharacterService(&stubCharRepo{mainChars: maleMCs}, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})
		if _, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "Yun Hai", Type: "Hero", Gender: "Male", Systems: []string{"Universe A Cultivation"},
		}); !errors.Is(err, character.ErrRoleGenderMismatch) {
			t.Fatalf("Hero with only male MCs: err = %v, want %v", err, character.ErrRoleGenderMismatch)
		}
	})

	t.Run("fails for an unknown system and does not persist", func(t *testing.T) {
		chars := &stubCharRepo{}
		svc := service.NewCharacterService(chars, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})
		_, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "X", Type: "SideCharacter", Gender: "Male", Systems: []string{"Nope"},
		})
		if !errors.Is(err, port.ErrPowerSystemNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemNotFound)
		}
		if len(chars.saved) != 0 {
			t.Errorf("persisted despite invalid system: %+v", chars.saved)
		}
	})

	t.Run("propagates a duplicate-name error", func(t *testing.T) {
		chars := &stubCharRepo{saveErr: port.ErrCharacterExists}
		svc := service.NewCharacterService(chars, systems, &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})
		_, err := svc.CreateCharacter(context.Background(), port.CreateCharacterInput{
			Name: "Lin Feng", Type: "MainCharacter", Gender: "Female", Systems: []string{"Universe A Cultivation"},
		})
		if !errors.Is(err, port.ErrCharacterExists) {
			t.Fatalf("err = %v, want %v", err, port.ErrCharacterExists)
		}
	})
}

func TestCharacterService_MainCharacter(t *testing.T) {
	mcs := []character.Character{{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female}}
	svc := service.NewCharacterService(&stubCharRepo{mainChars: mcs}, knownSystems(), &stubSpeciesRepo{}, stubClassRepo{}, stubProfessionRepo{}, stubItemRepo{}, &stubWorld{}, service.CharacterDefaults{})

	got, err := svc.MainCharacter(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Lin Feng" {
		t.Errorf("Name = %q, want %q", got.Name, "Lin Feng")
	}
}

// stubWorld is a port.DefaultWorldProvisioner double that records that it was
// called and returns the Human base and default power-system name.
type stubWorld struct{ called bool }

func (s *stubWorld) EnsureDefaults(context.Context) (port.DefaultWorld, error) {
	s.called = true
	return port.DefaultWorld{
		PowerSystem: service.DefaultPowerSystem,
		Species:     character.HumanBase(),
	}, nil
}

type stubSpeciesRepo struct{}

func (s *stubSpeciesRepo) Save(ctx context.Context, sp character.Species) error { return nil }
func (s *stubSpeciesRepo) FindByName(ctx context.Context, name string) (character.Species, error) {
	return character.Species{Name: name}, nil
}
func (s *stubSpeciesRepo) List(ctx context.Context) ([]character.Species, error) { return nil, nil }
