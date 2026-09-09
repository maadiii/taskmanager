package task

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/errors"
)

func (r *repo) GetTaskByIdAndOwner(ctx context.Context, id, ownerId string) (*task.Entity, error) {
	ctx, span := r.tracer.Start(ctx, "repo.GetTaskByIdAndOwner")
	defer span.End()

	row := r.client.QueryRow(ctx, GetByIdAndOwnerQuery, id, ownerId)

	entity := new(task.Entity)
	err := row.Scan(
		&entity.ID,
		&entity.UserID,
		&entity.Title,
		&entity.Description,
		&entity.Status,
		&entity.Priority,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if nil == err {
		return entity, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.NotFound(err, "task_id")
	}

	return nil, err
}

const GetByIdAndOwnerQuery = `
SELECT
	id, 
	user_id,
	title, 
	description,
	status,
	priority,
	created_at,
	updated_at
FROM tasks WHERE id = $1 AND user_id = $2
`
