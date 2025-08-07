package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskPatch struct {
	Title       *string
	Description *string
	Status      *string
}

var ErrNotFound = errors.New("not found")

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	List(ctx context.Context, status string, limit, offset int) ([]Task, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id int64) error
}

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(p *pgxpool.Pool) *TaskRepo { return &TaskRepo{pool: p} }

func (r *TaskRepo) Create(ctx context.Context, t *Task) error {
	const q = `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, q, t.Title, t.Description, t.Status).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepo) List(ctx context.Context, status string, limit, offset int) ([]Task, error) {
	const q = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE ($1 = '' OR status = $1)
		ORDER BY id
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (*Task, error) {
	const q = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1`
	var t Task
	err := r.pool.QueryRow(ctx, q, id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepo) Update(ctx context.Context, t *Task) error {
	const q = `
		UPDATE tasks
		SET title=$2,
		    description=$3,
		    status=$4,
		    updated_at=$5
		WHERE id=$1`
	tag, err := r.pool.Exec(ctx, q, t.ID, t.Title, t.Description, t.Status, t.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM tasks WHERE id=$1`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
