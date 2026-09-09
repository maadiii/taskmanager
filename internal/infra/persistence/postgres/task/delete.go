package task

import (
	"context"

	"github.com/maadiii/taskmanager/pkg/errors"
)

func (r *repo) DeleteByIdAndOwner(ctx context.Context, id, ownerId string) error {
	_, err := r.client.Exec(ctx, DeleteByIdAndOwnerQuery, id, ownerId)
	if nil == err {
		return nil
	}

	return errors.Wrap(err)
}

const DeleteByIdAndOwnerQuery = `
DELETE FROM tasks
WHERE id = $1 AND user_id = $2
`
