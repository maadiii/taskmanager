package task

import (
	"context"

	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/utils"
)

func (r *repo) List(ctx context.Context, status, userId string, limit int, lastId string) ([]task.Entity, error) {
	ctx, span := r.tracer.Start(ctx, "repo.List")
	defer span.End()

	p := new(utils.Pagination).
		Select("tasks", "id", "title", "description", "status").
		Limit(limit).
		Desc().
		LastID(lastId)

	query, args := "", []any{} //nolint
	if status != "" {
		query, args = p.Paginate("status = $1 AND user_id = $2", status, userId)
	} else {
		query, args = p.Paginate("user_id = $1", userId)
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Entity
	for rows.Next() {
		var t task.Entity
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
