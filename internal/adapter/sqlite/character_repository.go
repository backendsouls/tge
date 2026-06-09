package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// CharacterRepository implements port.CharacterRepository over SQLite. Characters
// are keyed by name and linked to their power systems via character_systems. The
// systems are stored by name only (their trees live with the power systems).
type CharacterRepository struct {
	db *sql.DB
}

// NewCharacterRepository opens (creating if needed) the database at dsn.
func NewCharacterRepository(dsn string) (*CharacterRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &CharacterRepository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *CharacterRepository) Close() error {
	return r.db.Close()
}

// Save persists a character and its system links in a transaction, returning
// port.ErrCharacterExists if the name is taken.
func (r *CharacterRepository) Save(ctx context.Context, c character.Character) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	const insChar = `INSERT INTO characters
		(name, type, gender, species, power, age, lifespan, class, profession)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insChar,
		c.Name, string(c.Type), string(c.Gender), primarySpeciesName(c), c.PowerValue, c.Mortal.Age, c.Mortal.Lifespan,
		c.Class.Name, c.Profession.Name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return port.ErrCharacterExists
			}
		}
		return fmt.Errorf("save character: %w", err)
	}

	const insStats = `INSERT INTO character_stats (character, str, agi, intel, vit, dex, wis, cha, luk)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	st := c.Stats
	if _, err := tx.ExecContext(ctx, insStats, c.Name, st.STR, st.AGI, st.INT, st.VIT, st.DEX, st.WIS, st.CHA, st.LUK); err != nil {
		return fmt.Errorf("save character stats: %w", err)
	}

	const insSys = `INSERT INTO character_systems (character, system) VALUES (?, ?)`
	for _, ps := range c.PowerSystems {
		if _, err := tx.ExecContext(ctx, insSys, c.Name, ps.Name); err != nil {
			return fmt.Errorf("link system %q: %w", ps.Name, err)
		}
	}
	return tx.Commit()
}

// FindByName returns a character, or port.ErrCharacterNotFound.
func (r *CharacterRepository) FindByName(ctx context.Context, name string) (character.Character, error) {
	const q = `SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters WHERE name = ?`
	c, err := r.scanCharacter(r.db.QueryRowContext(ctx, q, name))
	if errors.Is(err, sql.ErrNoRows) {
		return character.Character{}, fmt.Errorf("%w: %q", port.ErrCharacterNotFound, name)
	}
	if err != nil {
		return character.Character{}, fmt.Errorf("find character: %w", err)
	}
	if err := r.loadSystems(ctx, &c); err != nil {
		return character.Character{}, err
	}
	if err := r.loadStats(ctx, &c); err != nil {
		return character.Character{}, err
	}
	if err := r.loadInventory(ctx, &c); err != nil {
		return character.Character{}, err
	}
	if err := r.loadCultivations(ctx, &c); err != nil {
		return character.Character{}, err
	}
	return c, nil
}

// AddItem adds quantity of an item to a character's inventory, merging with any
// existing stack of the same item.
func (r *CharacterRepository) AddItem(ctx context.Context, character, item string, quantity int) error {
	const q = `INSERT INTO character_items (character, item, quantity) VALUES (?, ?, ?)
		ON CONFLICT(character, item) DO UPDATE SET quantity = quantity + excluded.quantity`
	if _, err := r.db.ExecContext(ctx, q, character, item, quantity); err != nil {
		return fmt.Errorf("add item to inventory: %w", err)
	}
	return nil
}

// SaveCultivation upserts a character's cultivation state at one (system, path)
// node, replacing any existing state there.
func (r *CharacterRepository) SaveCultivation(ctx context.Context, character string, rec port.CultivationRecord) error {
	const q = `INSERT INTO character_cultivations
		(character, system, path, realm, level_number, level_name, points, progress)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(character, system, path) DO UPDATE SET
			realm = excluded.realm,
			level_number = excluded.level_number,
			level_name = excluded.level_name,
			points = excluded.points,
			progress = excluded.progress`
	_, err := r.db.ExecContext(ctx, q,
		character, rec.System, rec.Path, rec.Realm, rec.LevelNumber, rec.LevelName, rec.Points, rec.Progress)
	if err != nil {
		return fmt.Errorf("save cultivation: %w", err)
	}
	return nil
}

// MainCharacters returns every main character ordered by name.
func (r *CharacterRepository) MainCharacters(ctx context.Context) ([]character.Character, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters WHERE type = ? ORDER BY name`,
		string(character.MainCharacter))
	if err != nil {
		return nil, fmt.Errorf("query main characters: %w", err)
	}
	defer rows.Close()

	var chars []character.Character
	for rows.Next() {
		c, err := r.scanCharacter(rows)
		if err != nil {
			return nil, fmt.Errorf("scan main character: %w", err)
		}
		chars = append(chars, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate main characters: %w", err)
	}
	for i := range chars {
		if err := r.loadAssociations(ctx, &chars[i]); err != nil {
			return nil, err
		}
	}
	return chars, nil
}

// loadAssociations populates a character's systems, stats and inventory.
func (r *CharacterRepository) loadAssociations(ctx context.Context, c *character.Character) error {
	if err := r.loadSystems(ctx, c); err != nil {
		return err
	}
	if err := r.loadStats(ctx, c); err != nil {
		return err
	}
	if err := r.loadInventory(ctx, c); err != nil {
		return err
	}
	return r.loadCultivations(ctx, c)
}

// List returns all characters ordered by name, each with its systems.
func (r *CharacterRepository) List(ctx context.Context) ([]character.Character, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}
	defer rows.Close()

	var chars []character.Character
	for rows.Next() {
		c, err := r.scanCharacter(rows)
		if err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}
		chars = append(chars, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate characters: %w", err)
	}
	for i := range chars {
		if err := r.loadAssociations(ctx, &chars[i]); err != nil {
			return nil, err
		}
	}
	return chars, nil
}

func (r *CharacterRepository) scanCharacter(s scanner) (character.Character, error) {
	var (
		name, ctype, gender, species, power string
		age, lifespan                       int
		class, profession                   string
	)
	if err := s.Scan(&name, &ctype, &gender, &species, &power, &age, &lifespan, &class, &profession); err != nil {
		return character.Character{}, err
	}
	return character.Character{
		Name:       name,
		Type:       character.CharacterType(ctype),
		Gender:     character.Gender(gender),
		Species:    []character.Species{{Name: species}},
		PowerValue: power,
		Mortal:     character.Mortal{Age: age, Lifespan: lifespan},
		Class:      rpg.Class{Name: class},
		Profession: rpg.Profession{Name: profession},
	}, nil
}

// primarySpeciesName returns the name of a character's first species, or "" when
// it has none. The single species column stores only this primary species; any
// additional species a character gains are not persisted yet.
func primarySpeciesName(c character.Character) string {
	if len(c.Species) == 0 {
		return ""
	}
	return c.Species[0].Name
}

// loadStats populates c.Stats from the character_stats table. A character
// created before stats existed simply has the zero block.
func (r *CharacterRepository) loadStats(ctx context.Context, c *character.Character) error {
	var st rpg.Stats
	err := r.db.QueryRowContext(ctx,
		`SELECT str, agi, intel, vit, dex, wis, cha, luk FROM character_stats WHERE character = ?`, c.Name).
		Scan(&st.STR, &st.AGI, &st.INT, &st.VIT, &st.DEX, &st.WIS, &st.CHA, &st.LUK)
	if errors.Is(err, sql.ErrNoRows) {
		c.Stats = rpg.Stats{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("load character stats: %w", err)
	}
	c.Stats = st
	return nil
}

// loadInventory populates c.Inventory (item stacks) preserving insertion order.
func (r *CharacterRepository) loadInventory(ctx context.Context, c *character.Character) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT item, quantity FROM character_items WHERE character = ? ORDER BY rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character inventory: %w", err)
	}
	defer rows.Close()

	c.Inventory = rpg.Inventory{}
	for rows.Next() {
		var stack rpg.ItemStack
		if err := rows.Scan(&stack.Item, &stack.Quantity); err != nil {
			return fmt.Errorf("scan inventory item: %w", err)
		}
		c.Inventory.Items = append(c.Inventory.Items, stack)
	}
	return rows.Err()
}

// loadSystems populates c.PowerSystems (names only) preserving insertion order.
func (r *CharacterRepository) loadSystems(ctx context.Context, c *character.Character) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT system FROM character_systems WHERE character = ? ORDER BY rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character systems: %w", err)
	}
	defer rows.Close()

	c.PowerSystems = nil
	for rows.Next() {
		var system string
		if err := rows.Scan(&system); err != nil {
			return fmt.Errorf("scan character system: %w", err)
		}
		c.PowerSystems = append(c.PowerSystems, progression.PowerSystem{Name: system})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate character systems: %w", err)
	}
	return nil
}

// loadCultivations populates c.Power from character_cultivations: one PowerSystem
// instance per distinct system, each carrying its progressed path nodes (with a
// CultivationState). Rows are grouped preserving insertion order.
func (r *CharacterRepository) loadCultivations(ctx context.Context, c *character.Character) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT system, path, realm, level_number, level_name, points, progress
		 FROM character_cultivations WHERE character = ? ORDER BY rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character cultivations: %w", err)
	}
	defer rows.Close()

	c.Power = nil
	bySystem := map[string]int{} // system name -> index into c.Power
	for rows.Next() {
		var (
			system, path, realm, levelName string
			levelNumber, points            int
			progress                       float64
		)
		if err := rows.Scan(&system, &path, &realm, &levelNumber, &levelName, &points, &progress); err != nil {
			return fmt.Errorf("scan cultivation: %w", err)
		}
		node := progression.Power{
			Name: path,
			State: progression.CultivationState{
				Realm:    progression.Realm{Name: realm},
				Level:    progression.Level{Number: levelNumber, Name: levelName},
				Points:   points,
				Progress: progress,
			},
		}
		idx, ok := bySystem[system]
		if !ok {
			c.Power = append(c.Power, progression.PowerSystem{Name: system, Kind: progression.Cultivation})
			idx = len(c.Power) - 1
			bySystem[system] = idx
		}
		c.Power[idx].Powers = append(c.Power[idx].Powers, node)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate character cultivations: %w", err)
	}
	return nil
}

var _ port.CharacterRepository = (*CharacterRepository)(nil)
