package cosmology_test

import (
	"testing"
	"tge/internal/core/domain/cosmology"
)

func TestNewTimeline(t *testing.T) {
	if _, err := cosmology.NewTimeline("   "); err != cosmology.ErrInvalidTimelineName {
		t.Errorf("expected ErrInvalidTimelineName, got %v", err)
	}

	tl, err := cosmology.NewTimeline(" Prime ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl.Name != "Prime" {
		t.Errorf("expected name Prime, got %q", tl.Name)
	}
}

func TestTimeline_AddEvent(t *testing.T) {
	tl, _ := cosmology.NewTimeline("Prime")

	if err := tl.AddEvent(1, "   "); err != cosmology.ErrInvalidEventDescription {
		t.Errorf("expected ErrInvalidEventDescription, got %v", err)
	}

	// Insert out of order; events should be kept sorted by Order.
	if err := tl.AddEvent(2, "Heat death"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tl.AddEvent(1, " Big Bang "); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl.Events) != 2 || tl.Events[0].Order != 1 || tl.Events[0].Description != "Big Bang" {
		t.Fatalf("events not sorted/trimmed as expected: %+v", tl.Events)
	}

	if err := tl.AddEvent(1, "duplicate"); err == nil {
		t.Fatal("expected error adding a duplicate order")
	}
}
