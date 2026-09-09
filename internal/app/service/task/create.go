package task

import (
	"context"

	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Create(ctx *appcontext.Context, rq *dto.CreateTaskRq) (*dto.CreateTaskRs, error) {
	c, span := s.tracer.Start(ctx, "service.Create")
	defer span.End()

	// Force domain service(business) rules
	task, err := task.Create(ctx, rq.Title, rq.Description)
	if err != nil {
		return nil, err
	}

	// Force application service rules
	if err := s.create(c, task); err != nil {
		return nil, err
	}

	return s.createTaskRs(task)
}

func (s *service) create(ctx context.Context, task *task.Entity) error {
	return s.uow.Do(ctx, func(ctx context.Context, repo port.RepoFactory) error {
		if err := repo.Tasks().CreateNew(ctx, task); err != nil {
			return err
		}

		// Do another uow related that may will be failed(just return err to rollback)
		// if err := raiseTaskEvent(); err != nil {
		// 	return err
		// }

		return nil
	})
}

func (s *service) createTaskRs(task *task.Entity) (*dto.CreateTaskRs, error) {
	return &dto.CreateTaskRs{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status.String(),
		Priority:         task.Priority.String(),
		CreatedAtUnixSec: task.CreatedAt.Unix(),
	}, nil
}
