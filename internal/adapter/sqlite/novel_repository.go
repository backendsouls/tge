package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/novel"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// NovelRepository implements port.NovelRepository over SQLite. Novels are keyed
// by title; main_character is unique so a character leads at most one novel.
type NovelRepository struct {
	db *sql.DB
}

// NewNovelRepository opens (creating if needed) the database at dsn.
func NewNovelRepository(dsn string) (*NovelRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &NovelRepository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *NovelRepository) Close() error {
	return r.db.Close()
}

// Save inserts a novel, mapping unique-constraint violations to a port error.
func (r *NovelRepository) Save(ctx context.Context, n novel.Novel) error {
	const q = `INSERT INTO novels (title, main_character) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, q, n.Title, n.MainCharacter)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%w: %q", port.ErrNovelExists, n.Title)
			}
		}
		return fmt.Errorf("save novel: %w", err)
	}
	return nil
}

// List returns all novels ordered by title.
func (r *NovelRepository) List(ctx context.Context) ([]novel.Novel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT title, main_character FROM novels ORDER BY title`)
	if err != nil {
		return nil, fmt.Errorf("list novels: %w", err)
	}
	defer rows.Close()

	var novels []novel.Novel
	for rows.Next() {
		var n novel.Novel
		if err := rows.Scan(&n.Title, &n.MainCharacter); err != nil {
			return nil, fmt.Errorf("scan novel: %w", err)
		}
		novels = append(novels, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate novels: %w", err)
	}
	return novels, nil
}

// FindByMainCharacter returns the novel led by the named character, or
// port.ErrNovelNotFound.
func (r *NovelRepository) FindByMainCharacter(ctx context.Context, mainCharacter string) (novel.Novel, error) {
	const q = `SELECT title, main_character FROM novels WHERE main_character = ?`
	var n novel.Novel
	err := r.db.QueryRowContext(ctx, q, mainCharacter).Scan(&n.Title, &n.MainCharacter)
	if errors.Is(err, sql.ErrNoRows) {
		return novel.Novel{}, port.ErrNovelNotFound
	}
	if err != nil {
		return novel.Novel{}, fmt.Errorf("find novel by main character: %w", err)
	}
	return n, nil
}

// FindByTitle returns a novel with its volumes and chapters, or
// port.ErrNovelNotFound.
func (r *NovelRepository) FindByTitle(ctx context.Context, title string) (novel.Novel, error) {
	const q = `SELECT title, main_character FROM novels WHERE title = ?`
	var n novel.Novel
	err := r.db.QueryRowContext(ctx, q, title).Scan(&n.Title, &n.MainCharacter)
	if errors.Is(err, sql.ErrNoRows) {
		return novel.Novel{}, port.ErrNovelNotFound
	}
	if err != nil {
		return novel.Novel{}, fmt.Errorf("find novel: %w", err)
	}
	if err := r.loadStructure(ctx, &n); err != nil {
		return novel.Novel{}, err
	}
	return n, nil
}

// SaveStructure replaces the stored volumes and chapters of a novel in a single
// transaction.
func (r *NovelRepository) SaveStructure(ctx context.Context, n novel.Novel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM chapters WHERE novel = ?`, n.Title); err != nil {
		return fmt.Errorf("clear chapters: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM volumes WHERE novel = ?`, n.Title); err != nil {
		return fmt.Errorf("clear volumes: %w", err)
	}

	for _, v := range n.Volumes {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO volumes (novel, number, title) VALUES (?, ?, ?)`, n.Title, v.Number, v.Title); err != nil {
			return fmt.Errorf("insert volume %d: %w", v.Number, err)
		}
		for _, c := range v.Chapters {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO chapters (novel, volume, number, title) VALUES (?, ?, ?, ?)`,
				n.Title, v.Number, c.Number, c.Title); err != nil {
				return fmt.Errorf("insert chapter %d.%d: %w", v.Number, c.Number, err)
			}
		}
	}
	return tx.Commit()
}

// loadStructure populates a novel's volumes (ordered by number) and their
// chapters (ordered by number).
func (r *NovelRepository) loadStructure(ctx context.Context, n *novel.Novel) error {
	vrows, err := r.db.QueryContext(ctx,
		`SELECT number, title FROM volumes WHERE novel = ? ORDER BY number`, n.Title)
	if err != nil {
		return fmt.Errorf("load volumes: %w", err)
	}
	defer vrows.Close()

	byNumber := map[int]*novel.Volume{}
	n.Volumes = nil
	for vrows.Next() {
		var v novel.Volume
		if err := vrows.Scan(&v.Number, &v.Title); err != nil {
			return fmt.Errorf("scan volume: %w", err)
		}
		n.Volumes = append(n.Volumes, v)
	}
	if err := vrows.Err(); err != nil {
		return fmt.Errorf("iterate volumes: %w", err)
	}
	for i := range n.Volumes {
		byNumber[n.Volumes[i].Number] = &n.Volumes[i]
	}

	crows, err := r.db.QueryContext(ctx,
		`SELECT volume, number, title FROM chapters WHERE novel = ? ORDER BY volume, number`, n.Title)
	if err != nil {
		return fmt.Errorf("load chapters: %w", err)
	}
	defer crows.Close()

	for crows.Next() {
		var volume int
		var c novel.Chapter
		if err := crows.Scan(&volume, &c.Number, &c.Title); err != nil {
			return fmt.Errorf("scan chapter: %w", err)
		}
		if v := byNumber[volume]; v != nil {
			v.Chapters = append(v.Chapters, c)
		}
	}
	if err := crows.Err(); err != nil {
		return fmt.Errorf("iterate chapters: %w", err)
	}
	return nil
}

var _ port.NovelRepository = (*NovelRepository)(nil)
