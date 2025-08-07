package service

import (
	"SkillsRock_test/internal/repo"
	"context"
	"errors"
)

var ErrInvalidStatus = errors.New("invalid status")

type TaskService struct{ r *repo.TaskRepo }

func New(r *repo.TaskRepo) *TaskService { return &TaskService{r} }

func (s *TaskService) Create(ctx context.Context, t *repo.Task) error {
	if t.Title == "" {
		return errors.New("title required")
	}
	return s.r.Create(ctx, t)
}

/* List / Update / Delete → просто прокидывают вниз + бизнес-валидация */
