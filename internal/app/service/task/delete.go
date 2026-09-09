package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Delete(ctx *appcontext.Context, rq *dto.DeleteTaskRq) (*dto.DeleteTaskRs, error) {
	entity, err := s.repo.GetTaskByIdAndOwner(ctx, rq.ID, ctx.Identity.UserID)
	if err != nil {
		return nil, err
	}

	if err := entity.Delete(ctx); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteByIdAndOwner(ctx, entity.ID, entity.UserID); err != nil {
		return nil, err
	}

	return &dto.DeleteTaskRs{ID: entity.ID, Deleted: true}, nil
}
