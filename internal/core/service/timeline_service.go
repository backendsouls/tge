package service

import (
	"context"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

// TimelineService implements port.TimelineService over a TimelineRepository.
type TimelineService struct {
	timelines port.TimelineRepository
}

func NewTimelineService(timelines port.TimelineRepository) *TimelineService {
	return &TimelineService{timelines: timelines}
}

func (s *TimelineService) GetTimeline(ctx context.Context, owner port.LocationRef) (cosmology.Timeline, error) {
	return s.timelines.Find(ctx, owner)
}

func (s *TimelineService) AddEvent(ctx context.Context, in port.AddTimelineEventInput) (cosmology.Timeline, error) {
	t, err := s.timelines.Find(ctx, in.Owner)
	if err != nil {
		return cosmology.Timeline{}, err
	}
	if err := t.AddEvent(in.Order, in.Description); err != nil {
		return cosmology.Timeline{}, err
	}
	if err := s.timelines.AddEvent(ctx, in.Owner, cosmology.Event{Order: in.Order, Description: in.Description}); err != nil {
		return cosmology.Timeline{}, err
	}
	return t, nil
}

// defaultTimelineName derives the name of a location's auto-provisioned timeline.
func defaultTimelineName(locationName string) string {
	return locationName + " Timeline"
}

// ensureTimeline provisions the single default-named timeline a location owns,
// ignoring the case where it already exists (so it is safe to call on every
// creation).
func ensureTimeline(ctx context.Context, repo port.TimelineRepository, owner port.LocationRef) error {
	t, err := cosmology.NewTimeline(defaultTimelineName(owner.Name))
	if err != nil {
		return err
	}
	return ignoreExists(repo.Save(ctx, owner, t), port.ErrTimelineExists)
}
