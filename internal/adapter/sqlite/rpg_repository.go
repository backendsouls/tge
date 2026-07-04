package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// isUniqueConstraint reports whether err is a SQLite primary-key/unique
// violation, used to map storage conflicts onto domain "already exists" errors.
func isUniqueConstraint(err error) bool {
	if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
		return serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}

// --- Ability ---

type AbilityRepository struct{ db *sql.DB }

func NewAbilityRepository(dsn string) (*AbilityRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &AbilityRepository{db: db}, nil
}
func (r *AbilityRepository) Close() error { return r.db.Close() }

func (r *AbilityRepository) Save(ctx context.Context, item rpg.Ability) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO abilities (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrAbilityExists
		}
		return fmt.Errorf("save ability: %w", err)
	}
	return nil
}
func (r *AbilityRepository) FindByName(ctx context.Context, name string) (rpg.Ability, error) {
	var item rpg.Ability
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM abilities WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Ability{}, fmt.Errorf("%w: %q", port.ErrAbilityNotFound, name)
	}
	if err != nil {
		return rpg.Ability{}, fmt.Errorf("find ability: %w", err)
	}
	return item, nil
}
func (r *AbilityRepository) List(ctx context.Context) ([]rpg.Ability, error) {
	return queryList(ctx, r.db, `SELECT name, description, grade FROM abilities ORDER BY name`, func(rows *sql.Rows) (rpg.Ability, error) {
		var item rpg.Ability
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	})
}

// --- Skill ---

type SkillRepository struct{ db *sql.DB }

func NewSkillRepository(dsn string) (*SkillRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &SkillRepository{db: db}, nil
}
func (r *SkillRepository) Close() error { return r.db.Close() }

func (r *SkillRepository) Save(ctx context.Context, item rpg.Skill) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO skills (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrSkillExists
		}
		return fmt.Errorf("save skill: %w", err)
	}
	return nil
}
func (r *SkillRepository) FindByName(ctx context.Context, name string) (rpg.Skill, error) {
	var item rpg.Skill
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM skills WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Skill{}, fmt.Errorf("%w: %q", port.ErrSkillNotFound, name)
	}
	if err != nil {
		return rpg.Skill{}, fmt.Errorf("find skill: %w", err)
	}
	return item, nil
}
func (r *SkillRepository) List(ctx context.Context) ([]rpg.Skill, error) {
	return queryList(ctx, r.db, `SELECT name, description, grade FROM skills ORDER BY name`, func(rows *sql.Rows) (rpg.Skill, error) {
		var item rpg.Skill
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	})
}

// --- Item ---

type ItemRepository struct{ db *sql.DB }

func NewItemRepository(dsn string) (*ItemRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &ItemRepository{db: db}, nil
}
func (r *ItemRepository) Close() error { return r.db.Close() }

func (r *ItemRepository) Save(ctx context.Context, item rpg.Item) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO items (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrItemExists
		}
		return fmt.Errorf("save item: %w", err)
	}
	return nil
}
func (r *ItemRepository) FindByName(ctx context.Context, name string) (rpg.Item, error) {
	var item rpg.Item
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM items WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Item{}, fmt.Errorf("%w: %q", port.ErrItemNotFound, name)
	}
	if err != nil {
		return rpg.Item{}, fmt.Errorf("find item: %w", err)
	}
	return item, nil
}
func (r *ItemRepository) List(ctx context.Context) ([]rpg.Item, error) {
	return queryList(ctx, r.db, `SELECT name, description, grade FROM items ORDER BY name`, func(rows *sql.Rows) (rpg.Item, error) {
		var item rpg.Item
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	})
}

// --- Effect ---

type EffectRepository struct{ db *sql.DB }

func NewEffectRepository(dsn string) (*EffectRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &EffectRepository{db: db}, nil
}
func (r *EffectRepository) Close() error { return r.db.Close() }

func (r *EffectRepository) Save(ctx context.Context, e rpg.Effect) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO effects (name, kind, description) VALUES (?, ?, ?)`, e.Name, string(e.Kind), e.Description)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrEffectExists
		}
		return fmt.Errorf("save effect: %w", err)
	}
	return nil
}
func (r *EffectRepository) FindByName(ctx context.Context, name string) (rpg.Effect, error) {
	var e rpg.Effect
	var kind string
	err := r.db.QueryRowContext(ctx, `SELECT name, kind, description FROM effects WHERE name = ?`, name).Scan(&e.Name, &kind, &e.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Effect{}, fmt.Errorf("%w: %q", port.ErrEffectNotFound, name)
	}
	if err != nil {
		return rpg.Effect{}, fmt.Errorf("find effect: %w", err)
	}
	e.Kind = rpg.EffectKind(kind)
	return e, nil
}
func (r *EffectRepository) List(ctx context.Context) ([]rpg.Effect, error) {
	return queryList(ctx, r.db, `SELECT name, kind, description FROM effects ORDER BY name`, func(rows *sql.Rows) (rpg.Effect, error) {
		var e rpg.Effect
		var kind string
		err := rows.Scan(&e.Name, &kind, &e.Description)
		e.Kind = rpg.EffectKind(kind)
		return e, err
	})
}

// --- Equipment ---

type EquipmentRepository struct{ db *sql.DB }

func NewEquipmentRepository(dsn string) (*EquipmentRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &EquipmentRepository{db: db}, nil
}
func (r *EquipmentRepository) Close() error { return r.db.Close() }

func (r *EquipmentRepository) Save(ctx context.Context, e rpg.Equipment) error {
	const q = `INSERT INTO equipment (name, slot, str, agi, intel, vit, dex, wis, cha, luk) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	b := e.Bonus
	_, err := r.db.ExecContext(ctx, q, e.Name, string(e.Slot), b.STR, b.AGI, b.INT, b.VIT, b.DEX, b.WIS, b.CHA, b.LUK)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrEquipmentExists
		}
		return fmt.Errorf("save equipment: %w", err)
	}
	return nil
}
func (r *EquipmentRepository) FindByName(ctx context.Context, name string) (rpg.Equipment, error) {
	e, err := scanEquipment(r.db.QueryRowContext(ctx,
		`SELECT name, slot, str, agi, intel, vit, dex, wis, cha, luk FROM equipment WHERE name = ?`, name))
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Equipment{}, fmt.Errorf("%w: %q", port.ErrEquipmentNotFound, name)
	}
	if err != nil {
		return rpg.Equipment{}, fmt.Errorf("find equipment: %w", err)
	}
	return e, nil
}
func (r *EquipmentRepository) List(ctx context.Context) ([]rpg.Equipment, error) {
	return queryList(ctx, r.db, `SELECT name, slot, str, agi, intel, vit, dex, wis, cha, luk FROM equipment ORDER BY name`, func(rows *sql.Rows) (rpg.Equipment, error) {
		return scanEquipment(rows)
	})
}

func scanEquipment(s scanner) (rpg.Equipment, error) {
	var e rpg.Equipment
	var slot string
	b := &e.Bonus
	if err := s.Scan(&e.Name, &slot, &b.STR, &b.AGI, &b.INT, &b.VIT, &b.DEX, &b.WIS, &b.CHA, &b.LUK); err != nil {
		return rpg.Equipment{}, err
	}
	e.Slot = rpg.EquipmentSlot(slot)
	return e, nil
}

// --- Profession ---

type ProfessionRepository struct{ db *sql.DB }

func NewProfessionRepository(dsn string) (*ProfessionRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &ProfessionRepository{db: db}, nil
}
func (r *ProfessionRepository) Close() error { return r.db.Close() }

func (r *ProfessionRepository) Save(ctx context.Context, item rpg.Profession) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO professions (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrProfessionExists
		}
		return fmt.Errorf("save profession: %w", err)
	}
	return nil
}
func (r *ProfessionRepository) FindByName(ctx context.Context, name string) (rpg.Profession, error) {
	var item rpg.Profession
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM professions WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Profession{}, fmt.Errorf("%w: %q", port.ErrProfessionNotFound, name)
	}
	if err != nil {
		return rpg.Profession{}, fmt.Errorf("find profession: %w", err)
	}
	return item, nil
}
func (r *ProfessionRepository) List(ctx context.Context) ([]rpg.Profession, error) {
	return queryList(ctx, r.db, `SELECT name, description, grade FROM professions ORDER BY name`, func(rows *sql.Rows) (rpg.Profession, error) {
		var item rpg.Profession
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	})
}

// --- Class ---

type ClassRepository struct{ db *sql.DB }

func NewClassRepository(dsn string) (*ClassRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &ClassRepository{db: db}, nil
}
func (r *ClassRepository) Close() error { return r.db.Close() }

func (r *ClassRepository) Save(ctx context.Context, item rpg.Class) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO classes (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrClassExists
		}
		return fmt.Errorf("save class: %w", err)
	}
	return nil
}
func (r *ClassRepository) FindByName(ctx context.Context, name string) (rpg.Class, error) {
	var item rpg.Class
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM classes WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Class{}, fmt.Errorf("%w: %q", port.ErrClassNotFound, name)
	}
	if err != nil {
		return rpg.Class{}, fmt.Errorf("find class: %w", err)
	}
	return item, nil
}
func (r *ClassRepository) List(ctx context.Context) ([]rpg.Class, error) {
	return queryList(ctx, r.db, `SELECT name, description, grade FROM classes ORDER BY name`, func(rows *sql.Rows) (rpg.Class, error) {
		var item rpg.Class
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	})
}

// --- Quest ---

type QuestRepository struct{ db *sql.DB }

func NewQuestRepository(dsn string) (*QuestRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &QuestRepository{db: db}, nil
}
func (r *QuestRepository) Close() error { return r.db.Close() }

func (r *QuestRepository) Save(ctx context.Context, q rpg.Quest) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO quests (name, description) VALUES (?, ?)`, q.Name, q.Description)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrQuestExists
		}
		return fmt.Errorf("save quest: %w", err)
	}
	return nil
}
func (r *QuestRepository) AddObjective(ctx context.Context, quest string, o rpg.Objective) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO quest_objectives (quest, ord, description) VALUES (?, ?, ?)`, quest, o.Order, o.Description)
	if err != nil {
		if isUniqueConstraint(err) {
			return rpg.ErrObjectiveOrderExists
		}
		return fmt.Errorf("add quest objective: %w", err)
	}
	return nil
}
func (r *QuestRepository) FindByName(ctx context.Context, name string) (rpg.Quest, error) {
	var q rpg.Quest
	err := r.db.QueryRowContext(ctx, `SELECT name, description FROM quests WHERE name = ?`, name).Scan(&q.Name, &q.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Quest{}, fmt.Errorf("%w: %q", port.ErrQuestNotFound, name)
	}
	if err != nil {
		return rpg.Quest{}, fmt.Errorf("find quest: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT ord, description FROM quest_objectives WHERE quest = ? ORDER BY ord`, name)
	if err != nil {
		return rpg.Quest{}, fmt.Errorf("load quest objectives: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var o rpg.Objective
		if err := rows.Scan(&o.Order, &o.Description); err != nil {
			return rpg.Quest{}, err
		}
		q.Objectives = append(q.Objectives, o)
	}
	return q, rows.Err()
}
func (r *QuestRepository) List(ctx context.Context) ([]rpg.Quest, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name, description FROM quests ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list quests: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var list []rpg.Quest
	for rows.Next() {
		var q rpg.Quest
		if err := rows.Scan(&q.Name, &q.Description); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

// --- Recipe ---

type RecipeRepository struct{ db *sql.DB }

func NewRecipeRepository(dsn string) (*RecipeRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &RecipeRepository{db: db}, nil
}
func (r *RecipeRepository) Close() error { return r.db.Close() }

func (r *RecipeRepository) Save(ctx context.Context, rc rpg.Recipe) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO recipes (name, output) VALUES (?, ?)`, rc.Name, rc.Output)
	if err != nil {
		if isUniqueConstraint(err) {
			return port.ErrRecipeExists
		}
		return fmt.Errorf("save recipe: %w", err)
	}
	return nil
}
func (r *RecipeRepository) AddInput(ctx context.Context, recipe string, in rpg.Ingredient) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO recipe_inputs (recipe, item, quantity) VALUES (?, ?, ?)`, recipe, in.Item, in.Quantity)
	if err != nil {
		if isUniqueConstraint(err) {
			return rpg.ErrIngredientExists
		}
		return fmt.Errorf("add recipe input: %w", err)
	}
	return nil
}
func (r *RecipeRepository) FindByName(ctx context.Context, name string) (rpg.Recipe, error) {
	var rc rpg.Recipe
	err := r.db.QueryRowContext(ctx, `SELECT name, output FROM recipes WHERE name = ?`, name).Scan(&rc.Name, &rc.Output)
	if errors.Is(err, sql.ErrNoRows) {
		return rpg.Recipe{}, fmt.Errorf("%w: %q", port.ErrRecipeNotFound, name)
	}
	if err != nil {
		return rpg.Recipe{}, fmt.Errorf("find recipe: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT item, quantity FROM recipe_inputs WHERE recipe = ? ORDER BY rowid`, name)
	if err != nil {
		return rpg.Recipe{}, fmt.Errorf("load recipe inputs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var in rpg.Ingredient
		if err := rows.Scan(&in.Item, &in.Quantity); err != nil {
			return rpg.Recipe{}, err
		}
		rc.Inputs = append(rc.Inputs, in)
	}
	return rc, rows.Err()
}
func (r *RecipeRepository) List(ctx context.Context) ([]rpg.Recipe, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name, output FROM recipes ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list recipes: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var list []rpg.Recipe
	for rows.Next() {
		var rc rpg.Recipe
		if err := rows.Scan(&rc.Name, &rc.Output); err != nil {
			return nil, err
		}
		list = append(list, rc)
	}
	return list, rows.Err()
}

// Interface checks.
var (
	_ port.AbilityRepository    = (*AbilityRepository)(nil)
	_ port.SkillRepository      = (*SkillRepository)(nil)
	_ port.ItemRepository       = (*ItemRepository)(nil)
	_ port.EffectRepository     = (*EffectRepository)(nil)
	_ port.EquipmentRepository  = (*EquipmentRepository)(nil)
	_ port.ProfessionRepository = (*ProfessionRepository)(nil)
	_ port.ClassRepository      = (*ClassRepository)(nil)
	_ port.QuestRepository      = (*QuestRepository)(nil)
	_ port.RecipeRepository     = (*RecipeRepository)(nil)
)

func queryList[T any](ctx context.Context, db *sql.DB, query string, scan func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var list []T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
