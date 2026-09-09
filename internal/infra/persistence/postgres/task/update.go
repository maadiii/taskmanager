package task

import (
	"context"

	domain "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/errors"
)

func (r *repo) UpdateByIdAndOwner(ctx context.Context, entity *domain.Entity) error {
	_, err := r.client.Exec(
		ctx,
		UpdateByIdAndOwnerQuery,
		entity.ID,
		entity.UserID,
		entity.Title,
		entity.Description,
		entity.Status,
		entity.Priority,
		entity.UpdatedAt,
	)
	if err == nil {
		return nil
	}

	return errors.Wrap(err)
}

const UpdateByIdAndOwnerQuery = `
UPDATE tasks
SET
	title = $3,
	description = $4,
	status = $5,
	priority = $6,
	updated_at = $7
WHERE id = $1 AND user_id = $2
`
