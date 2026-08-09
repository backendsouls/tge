package service

import (
	"context"
	"fmt"
	"sync"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

type CharacterService struct {
	chars       port.CharacterRepository
	systems     port.PowerSystemRepository
	species     port.SpeciesRepository
	classes     port.ClassRepository
	professions port.ProfessionRepository
	items       port.ItemRepository
	world       port.DefaultWorldProvisioner
	defaults    CharacterDefaults
	idle        *IdleService
	mu          sync.Mutex
}

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

func NewCharacterService(chars port.CharacterRepository, systems port.PowerSystemRepository, species port.SpeciesRepository, classes port.ClassRepository, professions port.ProfessionRepository, items port.ItemRepository, world port.DefaultWorldProvisioner, defaults CharacterDefaults) *CharacterService {
	idle := NewIdleService(chars)
	return &CharacterService{chars: chars, systems: systems, species: species, classes: classes, professions: professions, items: items, world: world, defaults: defaults, idle: idle}
}

func (s *CharacterService) CreateCharacter(ctx context.Context, in port.CreateCharacterInput) (character.Character, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctype := character.CharacterType(in.Type)

	if ctype == character.MainCharacter {
		dw, err := s.world.EnsureDefaults(ctx)
		if err != nil {
			return character.Character{}, err
		}
		if in.Species == "" {
			in.Species = dw.Species.Name
		}
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

	if in.Gender == "" {
		if sp.DefaultGender != "" {
			in.Gender = string(sp.DefaultGender)
		} else {
			in.Gender = s.defaults.gender()
		}
	}

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

	baseName := in.Name
	name := baseName
	suffix := 1
	for {
		_, err := s.chars.FindByName(ctx, name)
		if err == port.ErrCharacterNotFound {
			break
		}
		if err != nil {
			return character.Character{}, err
		}
		name = fmt.Sprintf("%s (%d)", baseName, suffix)
		suffix++
	}
	in.Name = name

	c, err := character.NewMortalCharacter(character.CharacterConfig{
		Name:       in.Name,
		Type:       ctype,
		Gender:     character.Gender(in.Gender),
		Species:    sp,
		Class:      class,
		Profession: profession,
		Age:        s.defaults.Age,
		Stats:      s.defaults.Stats,
	})
	if err != nil {
		return character.Character{}, err
	}
	if err := s.chars.Save(ctx, c); err != nil {
		return character.Character{}, err
	}
	return c, nil
}

func (s *CharacterService) GiveItem(ctx context.Context, in port.GiveItemInput) (character.Character, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, err := s.chars.FindByName(ctx, in.Character)
	if err != nil {
		return character.Character{}, err
	}
	if _, err := s.items.FindByName(ctx, in.Item); err != nil {
		return character.Character{}, err
	}
	if in.Quantity <= 0 {
		return character.Character{}, rpg.ErrInvalidQuantity
	}

	// Add to inventory
	found := false
	for i, stack := range c.Inventory.Items {
		if stack.Item == in.Item {
			c.Inventory.Items[i].Quantity += in.Quantity
			found = true
			break
		}
	}
	if !found {
		c.Inventory.Items = append(c.Inventory.Items, rpg.ItemStack{Item: in.Item, Quantity: in.Quantity})
	}

	if err := s.chars.Save(ctx, c); err != nil {
		return character.Character{}, err
	}
	return c, nil
}

// TrainNode attempts to progress or unlock a node within a character's power system graph.
func (s *CharacterService) TrainNode(ctx context.Context, in port.TrainNodeInput) (character.Character, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, err := s.chars.FindByName(ctx, in.Character)
	if err != nil {
		return character.Character{}, err
	}

	sys, err := s.systems.FindByName(ctx, in.System)
	if err != nil {
		return character.Character{}, fmt.Errorf("train node: load system %q: %w", in.System, err)
	}

	node, ok := sys.Nodes[in.NodeID]
	if !ok {
		return character.Character{}, fmt.Errorf("train node: node %q does not exist in system %q", in.NodeID, in.System)
	}

	// Verify the character hasn't already unlocked this node (prevent duplicates)
	for _, progress := range c.UnlockedNodes {
		if progress.System == in.System && progress.NodeID == in.NodeID {
			return character.Character{}, fmt.Errorf("train node: character already unlocked node %q", in.NodeID)
		}
	}

	// Validation: Dependency Check (Parents)
	for _, parentID := range node.Parents {
		hasParent := false
		for _, p := range c.UnlockedNodes {
			if p.System == in.System && p.NodeID == parentID {
				hasParent = true
				break
			}
		}
		if !hasParent {
			return character.Character{}, fmt.Errorf("train node: missing required parent node %q", parentID)
		}
	}

	// Validation: Mutually Exclusive Check
	for _, exclID := range node.MutuallyExclusive {
		for _, p := range c.UnlockedNodes {
			if p.System == in.System && p.NodeID == exclID {
				return character.Character{}, fmt.Errorf("train node: cannot unlock %q because mutually exclusive node %q is active", in.NodeID, exclID)
			}
		}
	}

	// For MVP, we grant the node immediately with full level 1 progress
	c.UnlockedNodes = append(c.UnlockedNodes, character.NodeProgress{
		System:    in.System,
		NodeID:    in.NodeID,
		Level:     1,
		Progress:  100.0,
		BasePower: node.BasePower,
	})

	// Dynamically sync power
	c.CalculateTotalPower()

	if err := s.chars.Save(ctx, c); err != nil {
		return character.Character{}, err
	}

	return c, nil
}

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

func (s *CharacterService) Character(ctx context.Context, name string) (character.Character, error) {
	return s.chars.FindByName(ctx, name)
}

func (s *CharacterService) ListCharacters(ctx context.Context) ([]character.Character, error) {
	return s.chars.List(ctx)
}

func (s *CharacterService) CleanCharacters(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chars.Clean(ctx)
}

func (s *CharacterService) AssignIdleActivity(ctx context.Context, charName string, activity string) (character.Character, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.idle.AssignActivity(ctx, charName, activity)
}

func (s *CharacterService) PassTime(ctx context.Context, charName string, seconds int64) (character.Character, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, err := s.chars.FindByName(ctx, charName)
	if err != nil {
		return character.Character{}, err
	}

	// Make sure we commit gains up to the previous NovelTime before advancing time
	s.idle.CommitOfflineGains(&c)

	c.NovelTime += seconds

	if err := s.chars.Save(ctx, c); err != nil {
		return character.Character{}, err
	}

	return c, nil
}
