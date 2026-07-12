package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
	"tge/internal/logger"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// CharacterRepository implements port.CharacterRepository over SQLite. Characters
// are keyed by name and linked to their power systems via character_systems. The
// systems are stored by name only (their trees live with the power systems).
type CharacterRepository struct {
	db *sql.DB
}

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
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
	defer func() { _ = tx.Rollback() }()

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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return character.Character{}, fmt.Errorf("begin read tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	const q = `SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters WHERE name = ?`
	c, err := r.scanCharacter(tx.QueryRowContext(ctx, q, name))
	if errors.Is(err, sql.ErrNoRows) {
		return character.Character{}, fmt.Errorf("%w: %q", port.ErrCharacterNotFound, name)
	}
	if err != nil {
		return character.Character{}, fmt.Errorf("find character: %w", err)
	}
	if err := r.loadAssociations(ctx, tx, &c); err != nil {
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
		(character, system, path, realm, level_number, level_name,
		 breakthrough_points, bottleneck_points, points, bottleneck, progress)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(character, system, path) DO UPDATE SET
			realm = excluded.realm,
			level_number = excluded.level_number,
			level_name = excluded.level_name,
			breakthrough_points = excluded.breakthrough_points,
			bottleneck_points = excluded.bottleneck_points,
			points = excluded.points,
			bottleneck = excluded.bottleneck,
			progress = excluded.progress`
	_, err := r.db.ExecContext(ctx, q,
		character, rec.System, rec.Path, rec.Realm, rec.LevelNumber, rec.LevelName,
		rec.BreakthroughPoints, rec.BottleneckPoints, rec.Points, rec.Bottleneck, rec.Progress)
	if err != nil {
		return fmt.Errorf("save cultivation: %w", err)
	}
	return nil
}

// SaveSuperPower upserts a character's superpower state.
func (r *CharacterRepository) SaveSuperPower(ctx context.Context, character string, rec port.SuperPowerRecord) error {
	const q = `INSERT INTO character_superpowers (character, system, path, tier)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(character, system, path) DO UPDATE SET tier = excluded.tier`
	if _, err := r.db.ExecContext(ctx, q, character, rec.System, rec.Path, rec.Tier); err != nil {
		return fmt.Errorf("save superpower: %w", err)
	}
	return nil
}

// MainCharacters returns every main character ordered by name.
func (r *CharacterRepository) MainCharacters(ctx context.Context) ([]character.Character, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("begin read tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx,
		`SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters WHERE type = ? ORDER BY name`,
		string(character.MainCharacter))
	if err != nil {
		return nil, fmt.Errorf("query main characters: %w", err)
	}
	defer func() { _ = rows.Close() }()

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
		if err := r.loadAssociations(ctx, tx, &chars[i]); err != nil {
			return nil, err
		}
	}
	return chars, nil
}

// loadAssociations populates a character's systems, stats and inventory.
func (r *CharacterRepository) loadAssociations(ctx context.Context, q queryer, c *character.Character) error {
	if err := r.loadSystems(ctx, q, c); err != nil {
		return err
	}
	if err := r.loadStats(ctx, q, c); err != nil {
		return err
	}
	if err := r.loadInventory(ctx, q, c); err != nil {
		return err
	}
	if err := r.loadCultivations(ctx, q, c); err != nil {
		return err
	}
	return r.loadSuperPowers(ctx, q, c)
}

// List returns all characters ordered by name, each with its systems.
func (r *CharacterRepository) List(ctx context.Context) ([]character.Character, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx,
		`SELECT name, type, gender, species, power, age, lifespan, class, profession FROM characters ORDER BY rowid`)
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}
	defer func() { _ = rows.Close() }()

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
		if err := r.loadAssociations(ctx, tx, &chars[i]); err != nil {
			return nil, err
		}
	}
	return chars, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

// scanCharacter converts a row into a Character domain entity.
func (r *CharacterRepository) scanCharacter(s rowScanner) (character.Character, error) {
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
func (r *CharacterRepository) loadStats(ctx context.Context, q queryer, c *character.Character) error {
	var st rpg.Stats
	err := q.QueryRowContext(ctx,
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
func (r *CharacterRepository) loadInventory(ctx context.Context, q queryer, c *character.Character) error {
	rows, err := q.QueryContext(ctx,
		`SELECT item, quantity FROM character_items WHERE character = ? ORDER BY rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character inventory: %w", err)
	}
	defer func() { _ = rows.Close() }()

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
func (r *CharacterRepository) loadSystems(ctx context.Context, q queryer, c *character.Character) error {
	rows, err := q.QueryContext(ctx,
		`SELECT cs.system, ps.kind FROM character_systems cs JOIN power_systems ps ON cs.system = ps.name WHERE cs.character = ? ORDER BY cs.rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character systems: %w", err)
	}
	defer func() { _ = rows.Close() }()

	c.PowerSystems = nil
	for rows.Next() {
		var system, kind string
		if err := rows.Scan(&system, &kind); err != nil {
			return fmt.Errorf("scan character system: %w", err)
		}
		c.PowerSystems = append(c.PowerSystems, progression.PowerSystem{Name: system, PowerSystemType: progression.PowerSystemType(kind)})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate character systems: %w", err)
	}
	return nil
}

// loadCultivations populates c.Power from character_cultivations: one PowerSystem
// instance per distinct system, each carrying its progressed path nodes (with a
// CultivationState). Rows are grouped preserving insertion order.
func (r *CharacterRepository) loadCultivations(ctx context.Context, q queryer, c *character.Character) error {
	rows, err := q.QueryContext(ctx,
		`SELECT cc.system, cc.path, cc.realm, cc.level_number, cc.level_name,
		        cc.breakthrough_points, cc.bottleneck_points, cc.points, cc.bottleneck, cc.progress
		 FROM character_cultivations cc JOIN power_systems ps ON cc.system = ps.name WHERE cc.character = ? AND ps.kind = 'Cultivation' ORDER BY cc.rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character cultivations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	c.Power = nil
	bySystem := map[string]int{} // system name -> index into c.Power
	for rows.Next() {
		var (
			system, path, realm, levelName                                        string
			levelNumber, breakthroughPoints, bottleneckPoints, points, bottleneck int
			progress                                                              float64
		)
		if err := rows.Scan(&system, &path, &realm, &levelNumber, &levelName,
			&breakthroughPoints, &bottleneckPoints, &points, &bottleneck, &progress); err != nil {
			return fmt.Errorf("scan cultivation: %w", err)
		}
		node := progression.Power{
			Name: path,
			State: progression.CultivationState{
				Realm: progression.Realm{Name: realm},
				Level: progression.Level{
					Number:             levelNumber,
					Name:               levelName,
					BreakthroughPoints: breakthroughPoints,
					BottleneckPoints:   bottleneckPoints,
				},
				Points:     points,
				Bottleneck: bottleneck,
				Progress:   progress,
			},
		}
		idx, ok := bySystem[system]
		if !ok {
			c.Power = append(c.Power, progression.PowerSystem{Name: system, PowerSystemType: progression.Cultivation})
			idx = len(c.Power) - 1
			bySystem[system] = idx
		} else if c.Power[idx].PowerSystemType != progression.Cultivation {
			// Log a warning if DB corruption exists
			logger.Dev("warning: skipping corrupted mixed-state cultivation node %q in system %q", path, system)
			continue // safeguard against mixed-state corruption
		}
		c.Power[idx].Powers = append(c.Power[idx].Powers, node)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate character cultivations: %w", err)
	}
	return nil
}

// loadSuperPowers populates c.Power from character_superpowers. It merges with any existing systems loaded.
func (r *CharacterRepository) loadSuperPowers(ctx context.Context, q queryer, c *character.Character) error {
	rows, err := q.QueryContext(ctx,
		`SELECT cs.system, cs.path, cs.tier
		 FROM character_superpowers cs JOIN power_systems ps ON cs.system = ps.name WHERE cs.character = ? AND ps.kind = 'SuperPower' ORDER BY cs.rowid`, c.Name)
	if err != nil {
		return fmt.Errorf("load character superpowers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	bySystem := map[string]int{}
	for i, ps := range c.Power {
		bySystem[ps.Name] = i
	}

	for rows.Next() {
		var (
			system, path string
			tier         int
		)
		if err := rows.Scan(&system, &path, &tier); err != nil {
			return fmt.Errorf("scan superpower: %w", err)
		}
		node := progression.Power{
			Name: path,
			State: progression.SuperPowerState{
				Tier: tier,
			},
		}
		idx, ok := bySystem[system]
		if !ok {
			c.Power = append(c.Power, progression.PowerSystem{Name: system, PowerSystemType: progression.SuperPower})
			idx = len(c.Power) - 1
			bySystem[system] = idx
		} else if c.Power[idx].PowerSystemType != progression.SuperPower {
			logger.Dev("warning: skipping corrupted mixed-state superpower node %q in system %q", path, system)
			continue // safeguard against mixed-state corruption
		}
		c.Power[idx].Powers = append(c.Power[idx].Powers, node)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate character superpowers: %w", err)
	}
	return nil
}

// UpdatePowerValue updates the computed power string for a character.
func (r *CharacterRepository) UpdatePowerValue(ctx context.Context, name string, power string) error {
	const q = `UPDATE characters SET power = ? WHERE name = ?`
	if _, err := r.db.ExecContext(ctx, q, power, name); err != nil {
		return fmt.Errorf("update power value: %w", err)
	}
	return nil
}

var _ port.CharacterRepository = (*CharacterRepository)(nil)
