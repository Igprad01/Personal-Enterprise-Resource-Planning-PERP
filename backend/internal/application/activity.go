package application

import (
	"context"
	"strings"
	"time"

	"personal-erp-backend/internal/domain"
	"personal-erp-backend/internal/repository"
)

type ActivityService struct {
	activities *repository.ActivityRepository
	categories *repository.CategoryRepository
}

func NewActivityService(activities *repository.ActivityRepository, categories *repository.CategoryRepository) *ActivityService {
	return &ActivityService{activities: activities, categories: categories}
}

func (s *ActivityService) List(ctx context.Context, f domain.ActivityFilter) ([]domain.Activity, error) {
	return s.activities.List(ctx, f)
}

func (s *ActivityService) Get(ctx context.Context, id int64) (domain.Activity, error) {
	return s.activities.Get(ctx, id)
}

func (s *ActivityService) Create(ctx context.Context, in *domain.Activity) (domain.Activity, error) {
	if err := s.validate(ctx, in); err != nil {
		return domain.Activity{}, err
	}
	return s.activities.Create(ctx, in)
}

func (s *ActivityService) Update(ctx context.Context, id int64, in *domain.Activity) (domain.Activity, error) {
	if err := s.validate(ctx, in); err != nil {
		return domain.Activity{}, err
	}
	return s.activities.Update(ctx, id, in)
}

func (s *ActivityService) Delete(ctx context.Context, id int64) error {
	return s.activities.Delete(ctx, id)
}

func (s *ActivityService) Stats(ctx context.Context, start, end string) (domain.Stats, error) {
	return s.activities.Stats(ctx, start, end)
}

func (s *ActivityService) validate(ctx context.Context, in *domain.Activity) error {
	in.Title = strings.TrimSpace(in.Title)
	in.ActivityDate = strings.TrimSpace(in.ActivityDate)

	if in.Title == "" {
		return &domain.ValidationError{Message: "title is required"}
	}
	if len(in.Title) > 255 {
		return &domain.ValidationError{Message: "title must be at most 255 characters"}
	}
	if _, err := time.Parse("2006-01-02", in.ActivityDate); err != nil {
		return &domain.ValidationError{Message: "activity_date must be a valid date (YYYY-MM-DD)"}
	}

	status := strings.ToLower(strings.TrimSpace(in.Status))
	if status == "" {
		status = "planned"
	}
	if !validStatus(status) {
		return &domain.ValidationError{Message: "status must be planned, in_progress, done, or cancelled"}
	}
	in.Status = status

	priority := strings.ToLower(strings.TrimSpace(in.Priority))
	if priority == "" {
		priority = "medium"
	}
	if !validPriority(priority) {
		return &domain.ValidationError{Message: "priority must be low, medium, or high"}
	}
	in.Priority = priority

	if in.CategoryID != nil {
		exists, err := s.categories.Exists(ctx, *in.CategoryID)
		if err != nil {
			return &domain.ValidationError{Message: "category lookup failed"}
		}
		if !exists {
			return &domain.ValidationError{Message: "category_id does not exist"}
		}
	}

	var start, end *string
	if in.StartTime != nil && strings.TrimSpace(*in.StartTime) != "" {
		t := strings.TrimSpace(*in.StartTime)
		if !validTime(t) {
			return &domain.ValidationError{Message: "start_time must be HH:MM or HH:MM:SS"}
		}
		start = &t
	}
	if in.EndTime != nil && strings.TrimSpace(*in.EndTime) != "" {
		t := strings.TrimSpace(*in.EndTime)
		if !validTime(t) {
			return &domain.ValidationError{Message: "end_time must be HH:MM or HH:MM:SS"}
		}
		end = &t
	}
	in.StartTime = start
	in.EndTime = end

	duration := in.DurationMin
	if start != nil && end != nil {
		if *end < *start {
			return &domain.ValidationError{Message: "end_time must be after start_time"}
		}
		d := computeMinutes(*start, *end)
		duration = &d
	}
	if duration != nil && *duration < 0 {
		return &domain.ValidationError{Message: "duration_min cannot be negative"}
	}
	in.DurationMin = duration

	return nil
}

func validStatus(v string) bool {
	switch v {
	case "planned", "in_progress", "done", "cancelled":
		return true
	}
	return false
}

func validPriority(v string) bool {
	switch v {
	case "low", "medium", "high":
		return true
	}
	return false
}

func validTime(v string) bool {
	for _, layout := range []string{"15:04", "15:04:05"} {
		if _, err := time.Parse(layout, v); err == nil {
			return true
		}
	}
	return false
}

func computeMinutes(start, end string) int {
	layout := "15:04:05"
	parse := func(v string) time.Time {
		if len(v) == 5 {
			t, err := time.Parse("15:04", v)
			if err != nil {
				return time.Time{}
			}
			return t
		}
		t, _ := time.Parse(layout, v)
		return t
	}
	s := parse(start)
	e := parse(end)
	d := e.Sub(s)
	if d < 0 {
		return 0
	}
	return int(d.Minutes())
}