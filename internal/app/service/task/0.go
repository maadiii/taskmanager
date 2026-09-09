package task

import (
	"context"
	"log"

	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	repo   port.TaskRepo
	uow    uow.UoW[port.RepoFactory]
	cache  port.TaskCache
	tracer trace.Tracer
}

func NewService(
	repo port.TaskRepo,
	uow uow.UoW[port.RepoFactory],
	cache port.TaskCache,
	tracer trace.Tracer,
) *service {
	return &service{repo: repo, uow: uow, cache: cache, tracer: tracer}
}

func (s *service) getCachedTask(ctx context.Context, userID, taskID string) *dto.GetByIdRs {
	cached, err := s.cache.Get(ctx, userID, taskID)
	if err != nil {
		log.Printf("%v", errors.Cache(err, "get task"))

		return nil
	}

	return cached
}

func (s *service) cacheTask(ctx context.Context, userID, taskID string, response *dto.GetByIdRs) {
	if err := s.cache.Set(ctx, userID, taskID, response); err != nil {
		log.Printf("%v", errors.Cache(err, "set task"))
	}
}

func (s *service) invalidateTaskCache(ctx context.Context, userID, taskID string) {
	if err := s.cache.Delete(ctx, userID, taskID); err != nil {
		log.Printf("%v", errors.Cache(err, "invalidate task"))
	}
}
