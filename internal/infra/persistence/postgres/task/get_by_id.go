package task

import (
	"context"

	"github.com/maadiii/taskmanager/internal/domain/task"
)

func (r *repo) GetTaskByIdAndOwner(ctx context.Context, id, ownerId string) (*task.Entity, error) {
	row := r.client.QueryRow(ctx, GetByIdQuery, id)

	entity := new(task.Entity)
	err := row.Scan(
		&entity.ID,
		&entity.Title,
		&entity.Description,
		&entity.Status,
		&entity.Priority,
		&entity.DueDate,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

const GetByIdQuery = `
SELECT
	id, 
	title, 
	description
	status,
	priority,
	due_date
	created_at,
	updated_at
FROM tasks WHERE id = $1 AND user_id = $2
`
