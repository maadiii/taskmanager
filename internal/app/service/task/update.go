package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Update(ctx *appcontext.Context, rq *dto.UpdateTaskRq) (*dto.UpdateTaskRs, error) {
	entity, err := s.repo.GetTaskByIdAndOwner(ctx, rq.ID, ctx.Identity.UserID)
	if err != nil {
		return nil, err
	}

	if err := entity.Update(ctx, rq.Title, rq.Description, rq.Status, rq.Priority); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateByIdAndOwner(ctx, entity); err != nil {
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
