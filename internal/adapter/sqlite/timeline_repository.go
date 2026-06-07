package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// TimelineRepository implements port.TimelineRepository over SQLite. Each
// location owns exactly one timeline, keyed by (owner_kind, owner_key); a
// realm's key is scoped by its universe since realm names are unique only there.
type TimelineRepository struct {
	db *sql.DB
}

func NewTimelineRepository(dsn string) (*TimelineRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &TimelineRepository{db: db}, nil
}

func (r *TimelineRepository) Close() error {
	return r.db.Close()
}

// ownerKey builds the storage key for an owning location. A realm is scoped by
// its universe; the cosmic levels are keyed by their globally unique name.
func ownerKey(owner port.LocationRef) string {
	if owner.Kind == port.LocationRealm {
		return owner.Universe + "\x1f" + owner.Name
	}
	return owner.Name
}

func (r *TimelineRepository) Save(ctx context.Context, owner port.LocationRef, t cosmology.Timeline) error {
	const q = `INSERT INTO timelines (owner_kind, owner_key, name) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, string(owner.Kind), ownerKey(owner), t.Name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return port.ErrTimelineExists
			}
		}
		return fmt.Errorf("save timeline: %w", err)
	}
	return nil
}

func (r *TimelineRepository) AddEvent(ctx context.Context, owner port.LocationRef, e cosmology.Event) error {
	const q = `INSERT INTO timeline_events (owner_kind, owner_key, ord, description) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, string(owner.Kind), ownerKey(owner), e.Order, e.Description)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return cosmology.ErrEventOrderExists
			}
		}
		return fmt.Errorf("add timeline event: %w", err)
	}
	return nil
}

func (r *TimelineRepository) Find(ctx context.Context, owner port.LocationRef) (cosmology.Timeline, error) {
	key := ownerKey(owner)
	var t cosmology.Timeline
	err := r.db.QueryRowContext(ctx,
		`SELECT name FROM timelines WHERE owner_kind = ? AND owner_key = ?`,
		string(owner.Kind), key).Scan(&t.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Timeline{}, fmt.Errorf("%w: %s %q", port.ErrTimelineNotFound, owner.Kind, owner.Name)
	}
	if err != nil {
		return cosmology.Timeline{}, fmt.Errorf("find timeline: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT ord, description FROM timeline_events WHERE owner_kind = ? AND owner_key = ? ORDER BY ord`,
		string(owner.Kind), key)
	if err != nil {
		return cosmology.Timeline{}, fmt.Errorf("load timeline events: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e cosmology.Event
		if err := rows.Scan(&e.Order, &e.Description); err != nil {
			return cosmology.Timeline{}, err
		}
		t.Events = append(t.Events, e)
	}
	return t, rows.Err()
}

var _ port.TimelineRepository = (*TimelineRepository)(nil)
