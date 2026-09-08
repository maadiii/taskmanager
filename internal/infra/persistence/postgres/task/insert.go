package task

import (
	"context"

	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/errors"
)

func (r *repo) CreateNew(ctx context.Context, entity *task.Entity) error {
	_, err := r.client.Exec(
		ctx, InsertQuery,
		entity.ID,
		entity.UserID,
		entity.Title,
		entity.Description,
		entity.Status,
		entity.Priority,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
	if nil == err {
		return nil
	}

	if errors.IsConstraint(err, "idx_unique_tasks_title_user_id") {
		return errors.UniqueViolation(err, "task", "title")
	}

	return err
}

const InsertQuery = `
INSERT INTO tasks(
	id, 
	user_id,
	title, 
	description,
	status,
	priority,
	created_at,
	updated_at
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`
