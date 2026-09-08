package task

import (
	"context"

	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Create(ctx *appcontext.Context, rq *dto.CreateTaskRq) (*dto.CreateTaskRs, error) {
	// Force domain service(business) rules
	task, err := task.Create(ctx, rq.Title, rq.Description)
	if err != nil {
		return nil, err
	}

	// Force application service rules
	if err := s.create(ctx, task); err != nil {
		return nil, err
	}

	return s.createTaskRs(task)
}

func (s *service) create(ctx *appcontext.Context, task *task.Entity) error {
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
		ID:       task.ID,
		Status:   task.Status.String(),
		Priority: task.Priority.String(),
	}, nil
}
