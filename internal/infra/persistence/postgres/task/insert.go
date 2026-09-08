package task

import (
	"context"
	"fmt"

	"github.com/maadiii/taskmanager/internal/domain/task"
)

func (r *repo) CreateNew(ctx context.Context, entity *task.Entity) error {
	_, err := r.client.Exec(
		ctx, InsertQuery,
		entity.ID,
		entity.Title,
		entity.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

const InsertQuery = `
INSERT INTO tasks(
	id, 
	title, 
	description
	status,
	priority,
	due_date
	created_at,
	updated_at
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`
