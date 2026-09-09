package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) Delete(ctx *appcontext.Context, rq *dto.DeleteTaskRq) (*dto.DeleteTaskRs, error) {
	c, span := s.tracer.Start(ctx, "service.Delete")
	defer span.End()

	entity, err := s.repo.GetTaskByIdAndOwner(c, rq.ID, ctx.Identity.UserID)
	if err != nil {
		return nil, err
	}

	if err := entity.Delete(ctx); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteByIdAndOwner(c, entity.ID, entity.UserID); err != nil {
		return nil, err
	}

	s.invalidateTaskCache(ctx, entity.UserID, entity.ID)

	return &dto.DeleteTaskRs{ID: entity.ID, Deleted: true}, nil
}
