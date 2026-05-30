package service_test

import (
	"context"
	"errors"
	"testing"

	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

// stubRepo is a hand-rolled test double for port.RealmRepository, letting us
// unit-test the service in isolation from any real storage.
type stubRepo struct {
	saved      []progression.Realm
	saveErr    error
	listResult []progression.Realm
	listErr    error
}

func (s *stubRepo) Save(_ context.Context, r progression.Realm) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved = append(s.saved, r)
	return nil
}

func (s *stubRepo) FindByName(context.Context, string) (progression.Realm, error) {
	return progression.Realm{}, port.ErrRealmNotFound
}

func (s *stubRepo) List(context.Context) ([]progression.Realm, error) {
	return s.listResult, s.listErr
}

func (s *stubRepo) AddLevel(context.Context, string, progression.Level) error { return nil }

func TestRealmService_AddRealm(t *testing.T) {
	t.Run("validates and persists a realm", func(t *testing.T) {
		repo := &stubRepo{}
		svc := service.NewRealmService(repo)

		got, err := svc.AddRealm(context.Background(), progression.RealmConfig{
			Name:            "Qi Condensation",
			PowerMultiplier: 2,
			PowerAdder:      10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "Qi Condensation" {
			t.Errorf("returned Name = %q, want %q", got.Name, "Qi Condensation")
		}
		if len(repo.saved) != 1 || repo.saved[0].Name != "Qi Condensation" {
			t.Fatalf("realm was not persisted: %+v", repo.saved)
		}
	})

	t.Run("rejects an invalid realm without touching the repository", func(t *testing.T) {
		repo := &stubRepo{}
		svc := service.NewRealmService(repo)

		_, err := svc.AddRealm(context.Background(), progression.RealmConfig{Name: ""})
		if !errors.Is(err, progression.ErrInvalidName) {
			t.Fatalf("err = %v, want %v", err, progression.ErrInvalidName)
		}
		if len(repo.saved) != 0 {
			t.Errorf("repository was written to on invalid input: %+v", repo.saved)
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		repo := &stubRepo{saveErr: port.ErrRealmExists}
		svc := service.NewRealmService(repo)

		_, err := svc.AddRealm(context.Background(), progression.RealmConfig{Name: "Foundation"})
		if !errors.Is(err, port.ErrRealmExists) {
			t.Fatalf("err = %v, want %v", err, port.ErrRealmExists)
		}
	})
}

func TestRealmService_ListRealms(t *testing.T) {
	want := []progression.Realm{{Name: "Qi Condensation"}, {Name: "Foundation"}}
	repo := &stubRepo{listResult: want}
	svc := service.NewRealmService(repo)

	got, err := svc.ListRealms(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d realms, want %d", len(got), len(want))
	}
}
