package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(p *pgxpool.Pool) *TaskRepo { return &TaskRepo{pool: p} }

func (r *TaskRepo) Create(ctx context.Context, t *Task) error {
	query := `INSERT INTO tasks (title, description) VALUES ($1,$2) RETURNING id,created_at,updated_at`
	return r.pool.QueryRow(ctx, query, t.Title, t.Description).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepo) List(ctx context.Context, status string, limit, offset int) ([]Task, error) {
	q := `SELECT id,title,description,status,created_at,updated_at
	      FROM tasks WHERE ($1 = '' OR status=$1)
	      ORDER BY id LIMIT $2 OFFSET $3`
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

/* Update, Delete аналогично */
