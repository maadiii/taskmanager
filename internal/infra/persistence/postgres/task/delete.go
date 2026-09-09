package task

import "context"

func (r *repo) DeleteByIdAndOwner(ctx context.Context, id, ownerId string) error {
	_, err := r.client.Exec(ctx, DeleteByIdAndOwnerQuery, id, ownerId)
	if err == nil {
		return nil
	}

	return err
}

const DeleteByIdAndOwnerQuery = `
DELETE FROM tasks
WHERE id = $1 AND user_id = $2
`
