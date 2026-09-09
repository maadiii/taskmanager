package task

import (
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func (s *service) List(ctx *appcontext.Context, rq *dto.ListTaskRq) (*dto.ListTaskRs, error) {
	entities, err := s.repo.List(ctx, rq.Status, ctx.Identity.UserID, rq.Limit, rq.LastID)
	if err != nil {
		return nil, err
	}

	response := &dto.ListTaskRs{
		Tasks: make([]dto.GetByIdRs, 0, len(entities)),
	}

	for i := range entities {
		response.Tasks = append(response.Tasks, *s.getByIdRs(&entities[i]))
	}

	return response, nil
}
