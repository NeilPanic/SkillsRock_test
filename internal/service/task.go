package service

import (
	"context"
	"errors"
	"github.com/NeilPanic/SkillsRock_test/internal/repo"
	"strings"
	"time"
)

var (
	ErrEmptyTitle      = errors.New("title is required")
	ErrInvalidStatus   = errors.New("invalid status value")
	ErrTaskNotFound    = errors.New("task not found")
	ErrStatusForbidden = errors.New("status transition not allowed")
)

var allowedStatus = map[string]struct{}{
	"new":         {},
	"in_progress": {},
	"done":        {},
}

var nextStatus = map[string][]string{
	"new":         {"in_progress", "done"},
	"in_progress": {"done"},
	"done":        {},
}

func isAllowedStatus(s string) bool {
	_, ok := allowedStatus[s]
	return ok
}

func canTransit(from, to string) bool {
	for _, n := range nextStatus[from] {
		if n == to {
			return true
		}
	}
	return false
}

type TaskService struct{ r *repo.TaskRepo }

func New(r *repo.TaskRepo) *TaskService { return &TaskService{r} }

func (s *TaskService) Create(ctx context.Context, t *repo.Task) error {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return ErrEmptyTitle
	}
	if t.Status == "" {
		t.Status = "new"
	}
	if !isAllowedStatus(t.Status) {
		return ErrInvalidStatus
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt
	return s.r.Create(ctx, t)
}

func (s *TaskService) List(ctx context.Context, status string, limit, offset int) ([]repo.Task, error) {
	if status != "" && !isAllowedStatus(status) {
		return nil, ErrInvalidStatus
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.r.List(ctx, status, limit, offset)
}

func (s *TaskService) Update(ctx context.Context, id int64, patch repo.TaskPatch) (*repo.Task, error) {
	current, err := s.r.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if patch.Title != nil {
		*patch.Title = strings.TrimSpace(*patch.Title)
		if *patch.Title == "" {
			return nil, ErrEmptyTitle
		}
		current.Title = *patch.Title
	}
	if patch.Description != nil {
		current.Description = *patch.Description
	}
	if patch.Status != nil {
		if !isAllowedStatus(*patch.Status) {
			return nil, ErrInvalidStatus
		}
		if !canTransit(current.Status, *patch.Status) {
			return nil, ErrStatusForbidden
		}
		current.Status = *patch.Status
	}

	current.UpdatedAt = time.Now()
	if err := s.r.Update(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	if err := s.r.Delete(ctx, id); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}
