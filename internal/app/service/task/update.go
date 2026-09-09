package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Update(ctx *appcontext.Context, rq *dto.UpdateTaskRq) (*dto.UpdateTaskRs, error) {
	c, span := s.tracer.Start(ctx, "service.Update")
	defer span.End()

	entity, err := s.repo.GetTaskByIdAndOwner(c, rq.ID, ctx.Identity.UserID)
	if err != nil {
		return nil, err
	}

	if err := entity.Update(ctx, rq.Title, rq.Description, rq.Status, rq.Priority); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateByIdAndOwner(c, entity); err != nil {
		return nil, err
	}

	s.invalidateTaskCache(ctx, entity.UserID, entity.ID)

	return &dto.UpdateTaskRs{
		ID:               entity.ID,
		Title:            entity.Title,
		Description:      entity.Description,
		Status:           entity.Status.String(),
		Priority:         entity.Priority.String(),
		CreatedAtUnixSec: entity.CreatedAt.Unix(),
		UpdatedAtUnixSec: entity.UpdatedAt.Unix(),
	}, nil
}
