package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) GetByID(ctx *appcontext.Context, rq *dto.GetByIdRq) (*dto.GetByIdRs, error) {
	if cached := s.getCachedTask(ctx, ctx.Identity.UserID, rq.ID); cached != nil {
		return cached, nil
	}

	task, err := s.repo.GetTaskByIdAndOwner(ctx, rq.ID, ctx.Identity.UserID)
	if err != nil {
		return nil, err
	}

	response := s.getByIdRs(task)
	s.cacheTask(ctx, ctx.Identity.UserID, rq.ID, response)

	return response, nil
}

func (s *service) getByIdRs(entity *task.Entity) *dto.GetByIdRs {
	return &dto.GetByIdRs{
		ID:               entity.ID,
		Title:            entity.Title,
		Description:      entity.Description,
		Status:           entity.Status.String(),
		Priority:         entity.Priority.String(),
		CreatedAtUnixSec: entity.CreatedAt.Unix(),
	}
}
