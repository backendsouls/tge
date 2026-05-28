package port

import (
	"context"
	"errors"

	"tge/internal/core/domain/rpg"
)

// Sentinel errors for the RPG entities.
var (
	ErrAbilityNotFound    = errors.New("ability: not found")
	ErrAbilityExists      = errors.New("ability: already exists")
	ErrSkillNotFound      = errors.New("skill: not found")
	ErrSkillExists        = errors.New("skill: already exists")
	ErrItemNotFound       = errors.New("item: not found")
	ErrItemExists         = errors.New("item: already exists")
	ErrEffectNotFound     = errors.New("effect: not found")
	ErrEffectExists       = errors.New("effect: already exists")
	ErrEquipmentNotFound  = errors.New("equipment: not found")
	ErrEquipmentExists    = errors.New("equipment: already exists")
	ErrProfessionNotFound = errors.New("profession: not found")
	ErrProfessionExists   = errors.New("profession: already exists")
	ErrClassNotFound      = errors.New("class: not found")
	ErrClassExists        = errors.New("class: already exists")
	ErrQuestNotFound      = errors.New("quest: not found")
	ErrQuestExists        = errors.New("quest: already exists")
	ErrRecipeNotFound     = errors.New("recipe: not found")
	ErrRecipeExists       = errors.New("recipe: already exists")
)

// Inputs.
type (
	CreateAbilityInput    struct{ Name, Description string }
	CreateSkillInput      struct{ Name, Description string }
	CreateItemInput       struct{ Name, Description string }
	CreateProfessionInput struct{ Name, Description string }
	CreateClassInput      struct{ Name, Description string }

	CreateEffectInput struct {
		Name        string
		Kind        string
		Description string
	}
	CreateEquipmentInput struct {
		Name  string
		Slot  string
		Bonus rpg.Stats
	}
	CreateQuestInput       struct{ Name, Description string }
	AddQuestObjectiveInput struct {
		Quest       string
		Order       int
		Description string
	}
	CreateRecipeInput  struct{ Name, Output string }
	AddIngredientInput struct {
		Recipe   string
		Item     string
		Quantity int
	}
)

// Driven ports (repositories).
type (
	AbilityRepository interface {
		Save(ctx context.Context, a rpg.Ability) error
		FindByName(ctx context.Context, name string) (rpg.Ability, error)
		List(ctx context.Context) ([]rpg.Ability, error)
	}
	SkillRepository interface {
		Save(ctx context.Context, s rpg.Skill) error
		FindByName(ctx context.Context, name string) (rpg.Skill, error)
		List(ctx context.Context) ([]rpg.Skill, error)
	}
	ItemRepository interface {
		Save(ctx context.Context, i rpg.Item) error
		FindByName(ctx context.Context, name string) (rpg.Item, error)
		List(ctx context.Context) ([]rpg.Item, error)
	}
	EffectRepository interface {
		Save(ctx context.Context, e rpg.Effect) error
		FindByName(ctx context.Context, name string) (rpg.Effect, error)
		List(ctx context.Context) ([]rpg.Effect, error)
	}
	EquipmentRepository interface {
		Save(ctx context.Context, e rpg.Equipment) error
		FindByName(ctx context.Context, name string) (rpg.Equipment, error)
		List(ctx context.Context) ([]rpg.Equipment, error)
	}
	ProfessionRepository interface {
		Save(ctx context.Context, p rpg.Profession) error
		FindByName(ctx context.Context, name string) (rpg.Profession, error)
		List(ctx context.Context) ([]rpg.Profession, error)
	}
	ClassRepository interface {
		Save(ctx context.Context, c rpg.Class) error
		FindByName(ctx context.Context, name string) (rpg.Class, error)
		List(ctx context.Context) ([]rpg.Class, error)
	}
	QuestRepository interface {
		Save(ctx context.Context, q rpg.Quest) error
		FindByName(ctx context.Context, name string) (rpg.Quest, error)
		List(ctx context.Context) ([]rpg.Quest, error)
		AddObjective(ctx context.Context, quest string, o rpg.Objective) error
	}
	RecipeRepository interface {
		Save(ctx context.Context, r rpg.Recipe) error
		FindByName(ctx context.Context, name string) (rpg.Recipe, error)
		List(ctx context.Context) ([]rpg.Recipe, error)
		AddInput(ctx context.Context, recipe string, in rpg.Ingredient) error
	}
)

// Driving ports (services).
type (
	AbilityService interface {
		CreateAbility(ctx context.Context, in CreateAbilityInput) (rpg.Ability, error)
		GetAbility(ctx context.Context, name string) (rpg.Ability, error)
		ListAbilities(ctx context.Context) ([]rpg.Ability, error)
	}
	SkillService interface {
		CreateSkill(ctx context.Context, in CreateSkillInput) (rpg.Skill, error)
		GetSkill(ctx context.Context, name string) (rpg.Skill, error)
		ListSkills(ctx context.Context) ([]rpg.Skill, error)
	}
	ItemService interface {
		CreateItem(ctx context.Context, in CreateItemInput) (rpg.Item, error)
		GetItem(ctx context.Context, name string) (rpg.Item, error)
		ListItems(ctx context.Context) ([]rpg.Item, error)
	}
	EffectService interface {
		CreateEffect(ctx context.Context, in CreateEffectInput) (rpg.Effect, error)
		GetEffect(ctx context.Context, name string) (rpg.Effect, error)
		ListEffects(ctx context.Context) ([]rpg.Effect, error)
	}
	EquipmentService interface {
		CreateEquipment(ctx context.Context, in CreateEquipmentInput) (rpg.Equipment, error)
		GetEquipment(ctx context.Context, name string) (rpg.Equipment, error)
		ListEquipment(ctx context.Context) ([]rpg.Equipment, error)
	}
	ProfessionService interface {
		CreateProfession(ctx context.Context, in CreateProfessionInput) (rpg.Profession, error)
		GetProfession(ctx context.Context, name string) (rpg.Profession, error)
		ListProfessions(ctx context.Context) ([]rpg.Profession, error)
	}
	ClassService interface {
		CreateClass(ctx context.Context, in CreateClassInput) (rpg.Class, error)
		GetClass(ctx context.Context, name string) (rpg.Class, error)
		ListClasses(ctx context.Context) ([]rpg.Class, error)
	}
	QuestService interface {
		CreateQuest(ctx context.Context, in CreateQuestInput) (rpg.Quest, error)
		AddObjective(ctx context.Context, in AddQuestObjectiveInput) (rpg.Quest, error)
		GetQuest(ctx context.Context, name string) (rpg.Quest, error)
		ListQuests(ctx context.Context) ([]rpg.Quest, error)
	}
	RecipeService interface {
		CreateRecipe(ctx context.Context, in CreateRecipeInput) (rpg.Recipe, error)
		AddIngredient(ctx context.Context, in AddIngredientInput) (rpg.Recipe, error)
		GetRecipe(ctx context.Context, name string) (rpg.Recipe, error)
		ListRecipes(ctx context.Context) ([]rpg.Recipe, error)
	}
)
