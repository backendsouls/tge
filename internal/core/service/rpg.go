package service

import (
	"context"

	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

// This file implements the driving ports for the RPG bounded context. Each
// service validates input through the domain constructor, then persists via its
// driven repository.

// --- Ability ---

type AbilityService struct{ repo port.AbilityRepository }

func NewAbilityService(repo port.AbilityRepository) *AbilityService {
	return &AbilityService{repo: repo}
}

func (s *AbilityService) CreateAbility(ctx context.Context, in port.CreateAbilityInput) (rpg.Ability, error) {
	a, err := rpg.NewAbility(in.Name, in.Description, in.Grade)
	if err != nil {
		return rpg.Ability{}, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return rpg.Ability{}, err
	}
	return a, nil
}
func (s *AbilityService) GetAbility(ctx context.Context, name string) (rpg.Ability, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *AbilityService) ListAbilities(ctx context.Context) ([]rpg.Ability, error) {
	return s.repo.List(ctx)
}

// --- Skill ---

type SkillService struct{ repo port.SkillRepository }

func NewSkillService(repo port.SkillRepository) *SkillService { return &SkillService{repo: repo} }

func (s *SkillService) CreateSkill(ctx context.Context, in port.CreateSkillInput) (rpg.Skill, error) {
	sk, err := rpg.NewSkill(in.Name, in.Description, in.Grade)
	if err != nil {
		return rpg.Skill{}, err
	}
	if err := s.repo.Save(ctx, sk); err != nil {
		return rpg.Skill{}, err
	}
	return sk, nil
}
func (s *SkillService) GetSkill(ctx context.Context, name string) (rpg.Skill, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *SkillService) ListSkills(ctx context.Context) ([]rpg.Skill, error) { return s.repo.List(ctx) }

// --- Item ---

type ItemService struct{ repo port.ItemRepository }

func NewItemService(repo port.ItemRepository) *ItemService { return &ItemService{repo: repo} }

func (s *ItemService) CreateItem(ctx context.Context, in port.CreateItemInput) (rpg.Item, error) {
	i, err := rpg.NewItem(in.Name, in.Description, in.Grade)
	if err != nil {
		return rpg.Item{}, err
	}
	if err := s.repo.Save(ctx, i); err != nil {
		return rpg.Item{}, err
	}
	return i, nil
}
func (s *ItemService) GetItem(ctx context.Context, name string) (rpg.Item, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *ItemService) ListItems(ctx context.Context) ([]rpg.Item, error) { return s.repo.List(ctx) }

// --- Effect ---

type EffectService struct{ repo port.EffectRepository }

func NewEffectService(repo port.EffectRepository) *EffectService { return &EffectService{repo: repo} }

func (s *EffectService) CreateEffect(ctx context.Context, in port.CreateEffectInput) (rpg.Effect, error) {
	e, err := rpg.NewEffect(in.Name, rpg.EffectKind(in.Kind), in.Description)
	if err != nil {
		return rpg.Effect{}, err
	}
	if err := s.repo.Save(ctx, e); err != nil {
		return rpg.Effect{}, err
	}
	return e, nil
}
func (s *EffectService) GetEffect(ctx context.Context, name string) (rpg.Effect, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *EffectService) ListEffects(ctx context.Context) ([]rpg.Effect, error) {
	return s.repo.List(ctx)
}

// --- Equipment ---

type EquipmentService struct{ repo port.EquipmentRepository }

func NewEquipmentService(repo port.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

func (s *EquipmentService) CreateEquipment(ctx context.Context, in port.CreateEquipmentInput) (rpg.Equipment, error) {
	e, err := rpg.NewEquipment(in.Name, rpg.EquipmentSlot(in.Slot), in.Bonus)
	if err != nil {
		return rpg.Equipment{}, err
	}
	if err := s.repo.Save(ctx, e); err != nil {
		return rpg.Equipment{}, err
	}
	return e, nil
}
func (s *EquipmentService) GetEquipment(ctx context.Context, name string) (rpg.Equipment, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *EquipmentService) ListEquipment(ctx context.Context) ([]rpg.Equipment, error) {
	return s.repo.List(ctx)
}

// --- Profession ---

type ProfessionService struct{ repo port.ProfessionRepository }

func NewProfessionService(repo port.ProfessionRepository) *ProfessionService {
	return &ProfessionService{repo: repo}
}

func (s *ProfessionService) CreateProfession(ctx context.Context, in port.CreateProfessionInput) (rpg.Profession, error) {
	p, err := rpg.NewProfession(in.Name, in.Description, in.Grade)
	if err != nil {
		return rpg.Profession{}, err
	}
	if err := s.repo.Save(ctx, p); err != nil {
		return rpg.Profession{}, err
	}
	return p, nil
}
func (s *ProfessionService) GetProfession(ctx context.Context, name string) (rpg.Profession, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *ProfessionService) ListProfessions(ctx context.Context) ([]rpg.Profession, error) {
	return s.repo.List(ctx)
}

// --- Class ---

type ClassService struct{ repo port.ClassRepository }

func NewClassService(repo port.ClassRepository) *ClassService { return &ClassService{repo: repo} }

func (s *ClassService) CreateClass(ctx context.Context, in port.CreateClassInput) (rpg.Class, error) {
	c, err := rpg.NewClass(in.Name, in.Description, in.Grade)
	if err != nil {
		return rpg.Class{}, err
	}
	if err := s.repo.Save(ctx, c); err != nil {
		return rpg.Class{}, err
	}
	return c, nil
}
func (s *ClassService) GetClass(ctx context.Context, name string) (rpg.Class, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *ClassService) ListClasses(ctx context.Context) ([]rpg.Class, error) { return s.repo.List(ctx) }

// --- Quest ---

type QuestService struct{ repo port.QuestRepository }

func NewQuestService(repo port.QuestRepository) *QuestService { return &QuestService{repo: repo} }

func (s *QuestService) CreateQuest(ctx context.Context, in port.CreateQuestInput) (rpg.Quest, error) {
	q, err := rpg.NewQuest(in.Name, in.Description)
	if err != nil {
		return rpg.Quest{}, err
	}
	if err := s.repo.Save(ctx, q); err != nil {
		return rpg.Quest{}, err
	}
	return q, nil
}
func (s *QuestService) AddObjective(ctx context.Context, in port.AddQuestObjectiveInput) (rpg.Quest, error) {
	q, err := s.repo.FindByName(ctx, in.Quest)
	if err != nil {
		return rpg.Quest{}, err
	}
	if err := q.AddObjective(in.Order, in.Description); err != nil {
		return rpg.Quest{}, err
	}
	if err := s.repo.AddObjective(ctx, in.Quest, rpg.Objective{Order: in.Order, Description: in.Description}); err != nil {
		return rpg.Quest{}, err
	}
	return q, nil
}
func (s *QuestService) GetQuest(ctx context.Context, name string) (rpg.Quest, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *QuestService) ListQuests(ctx context.Context) ([]rpg.Quest, error) { return s.repo.List(ctx) }

// --- Recipe ---

type RecipeService struct{ repo port.RecipeRepository }

func NewRecipeService(repo port.RecipeRepository) *RecipeService { return &RecipeService{repo: repo} }

func (s *RecipeService) CreateRecipe(ctx context.Context, in port.CreateRecipeInput) (rpg.Recipe, error) {
	r, err := rpg.NewRecipe(in.Name, in.Output)
	if err != nil {
		return rpg.Recipe{}, err
	}
	if err := s.repo.Save(ctx, r); err != nil {
		return rpg.Recipe{}, err
	}
	return r, nil
}
func (s *RecipeService) AddIngredient(ctx context.Context, in port.AddIngredientInput) (rpg.Recipe, error) {
	r, err := s.repo.FindByName(ctx, in.Recipe)
	if err != nil {
		return rpg.Recipe{}, err
	}
	if err := r.AddInput(in.Item, in.Quantity); err != nil {
		return rpg.Recipe{}, err
	}
	if err := s.repo.AddInput(ctx, in.Recipe, rpg.Ingredient{Item: in.Item, Quantity: in.Quantity}); err != nil {
		return rpg.Recipe{}, err
	}
	return r, nil
}
func (s *RecipeService) GetRecipe(ctx context.Context, name string) (rpg.Recipe, error) {
	return s.repo.FindByName(ctx, name)
}
func (s *RecipeService) ListRecipes(ctx context.Context) ([]rpg.Recipe, error) {
	return s.repo.List(ctx)
}
