package task

import (
	"github.com/maadiii/taskmanager/internal/application/dto"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *Service) Create(ctx *appcontext.Context, rq *dto.CreateTaskRq) (*dto.CreateTaskRs, error) {
	task := task.Create(rq.Title, rq.Description)

	err := s.repo.CreateNewTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return s.createTaskRs(task), nil
}

func (s *Service) createTaskRs(task *task.Entity) *dto.CreateTaskRs {
	return &dto.CreateTaskRs{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status.String(),
		Priority:    task.Priority.String(),
	}
}
