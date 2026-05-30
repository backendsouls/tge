// Package service implements the application's use cases (driving ports).
//
// It orchestrates the domain and depends only on driven ports (interfaces),
// never on concrete adapters — keeping the core decoupled from storage and IO.
package service

import (
	"context"

	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

// RealmService implements port.RealmService on top of a RealmRepository.
type RealmService struct {
	repo port.RealmRepository
}

// NewRealmService wires the service to a repository. It accepts the interface,
// not a concrete type, so any adapter can be substituted (LSP/DIP).
func NewRealmService(repo port.RealmRepository) *RealmService {
	return &RealmService{repo: repo}
}

// AddRealm validates the config via the domain and persists the resulting realm.
func (s *RealmService) AddRealm(ctx context.Context, cfg progression.RealmConfig) (progression.Realm, error) {
	r, err := progression.NewRealm(cfg)
	if err != nil {
		return progression.Realm{}, err
	}
	if err := s.repo.Save(ctx, r); err != nil {
		return progression.Realm{}, err
	}
	return r, nil
}

// ListRealms returns all persisted realms.
func (s *RealmService) ListRealms(ctx context.Context) ([]progression.Realm, error) {
	return s.repo.List(ctx)
}

// GetRealm returns a single realm (with its levels) by name.
func (s *RealmService) GetRealm(ctx context.Context, name string) (progression.Realm, error) {
	return s.repo.FindByName(ctx, name)
}

// AddLevel validates and adds an ordered level to an existing realm.
func (s *RealmService) AddLevel(ctx context.Context, in port.AddLevelInput) (progression.Realm, error) {
	r, err := s.repo.FindByName(ctx, in.Realm)
	if err != nil {
		return progression.Realm{}, err
	}
	if err := r.AddLevel(in.Number, in.Name, in.BreakthroughPoints, in.BottleneckPoints); err != nil {
		return progression.Realm{}, err
	}
	level := progression.Level{
		Number:             in.Number,
		Name:               in.Name,
		BreakthroughPoints: in.BreakthroughPoints,
		BottleneckPoints:   in.BottleneckPoints,
	}
	if err := s.repo.AddLevel(ctx, in.Realm, level); err != nil {
		return progression.Realm{}, err
	}
	return r, nil
}
