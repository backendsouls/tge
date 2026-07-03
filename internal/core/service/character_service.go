package service

import (
	"context"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

// CharacterService implements port.CharacterService. It creates mortal
// characters within existing power systems, enforcing the single-main-character
// and Hero/Heroine rules, and reports the main character and the roster.
type CharacterService struct {
	chars       port.CharacterRepository
	systems     port.PowerSystemRepository
	species     port.SpeciesRepository
	classes     port.ClassRepository
	professions port.ProfessionRepository
	items       port.ItemRepository
	world       port.DefaultWorldProvisioner
	defaults    CharacterDefaults
}

// CharacterDefaults holds the fill-in values applied to a new character when the
// caller does not supply them. A blank Gender falls back to Male, and the zero
// Age/Stats fall back to the domain defaults.
type CharacterDefaults struct {
	Gender string
	Age    int
	Stats  rpg.Stats
}

func (d CharacterDefaults) gender() string {
	if d.Gender == "" {
		return string(character.Male)
	}
	return d.Gender
}

// NewCharacterService wires the service to its repositories and the default-world
// provisioner used to back a name-only main character. The class, profession and
// item repositories validate a character's optional RPG attributes and inventory.
// defaults supplies the fill-in gender/age/stats for new characters.
func NewCharacterService(chars port.CharacterRepository, systems port.PowerSystemRepository, species port.SpeciesRepository, classes port.ClassRepository, professions port.ProfessionRepository, items port.ItemRepository, world port.DefaultWorldProvisioner, defaults CharacterDefaults) *CharacterService {
	return &CharacterService{chars: chars, systems: systems, species: species, classes: classes, professions: professions, items: items, world: world, defaults: defaults}
}

// CreateCharacter validates input, enforces cross-character role rules and
// persists a fresh mortal character.
func (s *CharacterService) CreateCharacter(ctx context.Context, in port.CreateCharacterInput) (character.Character, error) {
	ctype := character.CharacterType(in.Type)

	// A main character can be created from just a name: it is born from the
	// Human base into the default cosmology, which is provisioned in the
	// background. Explicitly provided values are kept.
	if ctype == character.MainCharacter {
		dw, err := s.world.EnsureDefaults(ctx)
		if err != nil {
			return character.Character{}, err
		}
		if in.Species == "" {
			in.Species = dw.Species.Name
		}
		if len(in.Systems) == 0 {
			in.Systems = []string{dw.PowerSystem}
		}
	}

	systems, err := s.resolveSystems(ctx, in.Systems)
	if err != nil {
		return character.Character{}, err
	}

	mcs, err := s.chars.MainCharacters(ctx)
	if err != nil {
		return character.Character{}, err
	}
	if err := character.CheckRole(ctype, mcs); err != nil {
		return character.Character{}, err
	}

	sp, err := s.species.FindByName(ctx, in.Species)
	if err != nil {
		return character.Character{}, err
	}

	// When no gender is given, fall back to the species' default gender (e.g.
	// Human defaults to Male), then to the global default.
	if in.Gender == "" {
		if sp.DefaultGender != "" {
			in.Gender = string(sp.DefaultGender)
		} else {
			in.Gender = s.defaults.gender()
		}
	}

	// Class and profession are optional, but must exist when given.
	var class rpg.Class
	if in.Class != "" {
		class, err = s.classes.FindByName(ctx, in.Class)
		if err != nil {
			return character.Character{}, err
		}
	}
	var profession rpg.Profession
	if in.Profession != "" {
		profession, err = s.professions.FindByName(ctx, in.Profession)
		if err != nil {
			return character.Character{}, err
		}
	}

	c, err := character.NewMortalCharacter(character.CharacterConfig{
		Name:         in.Name,
		Type:         ctype,
		Gender:       character.Gender(in.Gender),
		Species:      sp,
		PowerSystems: systems,
		Class:        class,
		Profession:   profession,
		Age:          s.defaults.Age,
		Stats:        s.defaults.Stats,
	})
	if err != nil {
		return character.Character{}, err
	}
	if err := s.chars.Save(ctx, c); err != nil {
		return character.Character{}, err
	}
	return c, nil
}

// GiveItem adds a quantity of an existing item to a character's inventory and
// returns the updated character.
func (s *CharacterService) GiveItem(ctx context.Context, in port.GiveItemInput) (character.Character, error) {
	if _, err := s.chars.FindByName(ctx, in.Character); err != nil {
		return character.Character{}, err
	}
	if _, err := s.items.FindByName(ctx, in.Item); err != nil {
		return character.Character{}, err
	}
	if in.Quantity <= 0 {
		return character.Character{}, rpg.ErrInvalidQuantity
	}
	if err := s.chars.AddItem(ctx, in.Character, in.Item, in.Quantity); err != nil {
		return character.Character{}, err
	}
	return s.chars.FindByName(ctx, in.Character)
}

// Cultivate sets a character's cultivation state at one power node and returns
// the updated character. System defaults to the character's first power system
// and Path to that system's name when not given.
func (s *CharacterService) Cultivate(ctx context.Context, in port.CultivateInput) (character.Character, error) {
	c, err := s.chars.FindByName(ctx, in.Character)
	if err != nil {
		return character.Character{}, err
	}

	system := in.System
	if system == "" {
		if len(c.PowerSystems) == 0 {
			return character.Character{}, character.ErrMissingSystem
		}
		system = c.PowerSystems[0].Name
	}
	path := in.Path
	if path == "" {
		path = system
	}

	rec := port.CultivationRecord{
		System:             system,
		Path:               path,
		Realm:              in.Realm,
		LevelNumber:        in.LevelNumber,
		LevelName:          in.LevelName,
		BreakthroughPoints: in.BreakthroughPoints,
		BottleneckPoints:   in.BottleneckPoints,
		Points:             in.Points,
		Bottleneck:         in.Bottleneck,
		Progress:           in.Progress,
	}
	if err := s.chars.SaveCultivation(ctx, in.Character, rec); err != nil {
		return character.Character{}, err
	}
	return s.chars.FindByName(ctx, in.Character)
}

// Train adds cultivation points to a character's power node, filling the current
// level's breakthrough gate then its bottleneck gate, breaking through to the
// next level and, at a realm's ceiling, into the next realm (by tier, from the
// caller-supplied ordered realms). Leftover points past the highest realm are
// discarded. System defaults to the character's first power system and Path to
// the System name.
func (s *CharacterService) Train(ctx context.Context, in port.TrainInput) (character.Character, error) {
	if in.Points <= 0 {
		return character.Character{}, port.ErrInvalidTrainingPoints
	}
	c, err := s.chars.FindByName(ctx, in.Character)
	if err != nil {
		return character.Character{}, err
	}

	system := in.System
	if system == "" {
		if len(c.PowerSystems) == 0 {
			return character.Character{}, character.ErrMissingSystem
		}
		system = c.PowerSystems[0].Name
	}
	path := in.Path
	if path == "" {
		path = system
	}

	// Establish the starting state, rehydrating the realm with its full level
	// list (persistence stores only the current realm/level).
	cs, ok := cultivationAt(c, system, path)
	if ok {
		full, found := realmByName(in.Realms, cs.Realm.Name)
		if !found {
			return character.Character{}, port.ErrRealmNotFound
		}
		cs.Realm = full
		if lvl, ok := levelOf(full, cs.Level.Number); ok {
			cs.Level = lvl
		}
	} else {
		if len(in.Realms) == 0 || len(in.Realms[0].Levels) == 0 {
			return character.Character{}, port.ErrNoRealms
		}
		cs = progression.CultivationState{Realm: in.Realms[0], Level: in.Realms[0].Levels[0]}
	}

	// Apply the points, crossing into the next realm by tier as gates fill.
	remaining := in.Points
	for {
		cs, remaining = cs.AdvanceWithin(remaining)
		if remaining == 0 {
			break
		}
		next, ok := realmAfter(in.Realms, cs.Realm.Name)
		if !ok || len(next.Levels) == 0 {
			break // highest realm reached; leftover points are discarded
		}
		cs.Realm = next
		cs.Level = next.Levels[0]
		cs.Points = 0
		cs.Bottleneck = 0
	}

	rec := port.CultivationRecord{
		System:             system,
		Path:               path,
		Realm:              cs.Realm.Name,
		LevelNumber:        cs.Level.Number,
		LevelName:          cs.Level.Name,
		BreakthroughPoints: cs.Level.BreakthroughPoints,
		BottleneckPoints:   cs.Level.BottleneckPoints,
		Points:             cs.Points,
		Bottleneck:         cs.Bottleneck,
		Progress:           cs.Progress,
	}
	if err := s.chars.SaveCultivation(ctx, in.Character, rec); err != nil {
		return character.Character{}, err
	}
	return s.chars.FindByName(ctx, in.Character)
}

// cultivationAt returns the CultivationState at a character's (system, path)
// node, if it has one.
func cultivationAt(c character.Character, system, path string) (progression.CultivationState, bool) {
	for _, ps := range c.Power {
		if ps.Name != system {
			continue
		}
		for _, p := range ps.Powers {
			if p.Name != path {
				continue
			}
			if cs, ok := p.State.(progression.CultivationState); ok {
				return cs, true
			}
		}
	}
	return progression.CultivationState{}, false
}

// realmByName finds a realm by name within an ordered realm slice.
func realmByName(realms []progression.Realm, name string) (progression.Realm, bool) {
	for _, r := range realms {
		if r.Name == name {
			return r, true
		}
	}
	return progression.Realm{}, false
}

// realmAfter returns the realm immediately following the named one in the
// tier-ordered slice.
func realmAfter(realms []progression.Realm, name string) (progression.Realm, bool) {
	for i, r := range realms {
		if r.Name == name && i+1 < len(realms) {
			return realms[i+1], true
		}
	}
	return progression.Realm{}, false
}

// levelOf finds a realm's level by number.
func levelOf(r progression.Realm, number int) (progression.Level, bool) {
	for _, l := range r.Levels {
		if l.Number == number {
			return l, true
		}
	}
	return progression.Level{}, false
}

// MainCharacter returns the first main character by name (the default status
// target), or port.ErrCharacterNotFound if there are none.
func (s *CharacterService) MainCharacter(ctx context.Context) (character.Character, error) {
	mcs, err := s.chars.MainCharacters(ctx)
	if err != nil {
		return character.Character{}, err
	}
	if len(mcs) == 0 {
		return character.Character{}, port.ErrCharacterNotFound
	}
	return mcs[0], nil
}

// Character returns a character by name.
func (s *CharacterService) Character(ctx context.Context, name string) (character.Character, error) {
	return s.chars.FindByName(ctx, name)
}

// ListCharacters returns every character.
func (s *CharacterService) ListCharacters(ctx context.Context) ([]character.Character, error) {
	return s.chars.List(ctx)
}

// resolveSystems loads each named power system, validating that it exists.
func (s *CharacterService) resolveSystems(ctx context.Context, names []string) ([]progression.PowerSystem, error) {
	systems := make([]progression.PowerSystem, 0, len(names))
	for _, name := range names {
		ps, err := s.systems.FindByName(ctx, name)
		if err != nil {
			return nil, err
		}
		systems = append(systems, ps)
	}
	return systems, nil
}
