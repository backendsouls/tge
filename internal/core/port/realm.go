// Package port declares the interfaces (ports) that connect the domain core to
// the outside world.
//
//   - Driven ports (e.g. RealmRepository) are implemented by outbound adapters
//     such as a SQLite or in-memory store.
//   - Driving ports (e.g. RealmService) are implemented by the application core
//     and consumed by inbound adapters such as the CLI.
//
// Both are owned by the core, so dependencies always point inward (DIP).
package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cultivation"
)

var (
	// ErrRealmNotFound is returned when a lookup finds no matching realm.
	ErrRealmNotFound = errors.New("realm: not found")
	// ErrRealmExists is returned when saving a realm whose name is already taken.
	ErrRealmExists = errors.New("realm: already exists")
)

// AddLevelInput describes a level to add to an existing realm.
type AddLevelInput struct {
	Realm              string
	Number             int
	Name               string
	BreakthroughPoints int
}

// RealmRepository is a driven port: it persists and retrieves realms. Adapters
// (SQLite, in-memory, ...) implement it so the core never depends on storage.
type RealmRepository interface {
	Save(ctx context.Context, r cultivation.Realm) error
	FindByName(ctx context.Context, name string) (cultivation.Realm, error)
	List(ctx context.Context) ([]cultivation.Realm, error)
	// AddLevel adds an ordered level to an existing realm.
	AddLevel(ctx context.Context, realm string, l cultivation.Level) error
}

// RealmService is a driving port: the use cases the application exposes. Inbound
// adapters (the CLI) depend on this interface, not on a concrete implementation.
type RealmService interface {
	AddRealm(ctx context.Context, cfg cultivation.RealmConfig) (cultivation.Realm, error)
	ListRealms(ctx context.Context) ([]cultivation.Realm, error)
	GetRealm(ctx context.Context, name string) (cultivation.Realm, error)
	AddLevel(ctx context.Context, in AddLevelInput) (cultivation.Realm, error)
}
