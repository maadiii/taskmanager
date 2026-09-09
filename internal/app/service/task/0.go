package task

import (
	"log"

	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/maadiii/taskmanager/pkg/errors"
)

type service struct {
	repo  port.TaskRepo
	uow   uow.UoW[port.RepoFactory]
	cache port.TaskCache
}

func NewService(
	repo port.TaskRepo,
	uow uow.UoW[port.RepoFactory],
	cache port.TaskCache,
) *service {
	return &service{repo: repo, uow: uow}
}

func (s *service) getCachedTask(ctx *appcontext.Context, userID, taskID string) *dto.GetByIdRs {
	cached, err := s.cache.Get(ctx, userID, taskID)
	if err != nil {
		log.Printf("%v", errors.Cache(err, "get task"))

		return nil
	}

	return cached
}

func (s *service) cacheTask(ctx *appcontext.Context, userID, taskID string, response *dto.GetByIdRs) {
	if err := s.cache.Set(ctx, userID, taskID, response); err != nil {
		log.Printf("%v", errors.Cache(err, "set task"))
	}
}

func (s *service) invalidateTaskCache(ctx *appcontext.Context, userID, taskID string) {
	if err := s.cache.Delete(ctx, userID, taskID); err != nil {
		log.Printf("%v", errors.Cache(err, "invalidate task"))
	}
}
