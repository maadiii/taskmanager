package port

import (
	"context"

	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

type TaskRepo interface {
	CreateNew(ctx context.Context, entity *task.Entity) error
	GetTaskByIdAndOwner(ctx context.Context, id, ownerId string) (*task.Entity, error)
}

type TaskService interface {
	Create(ctx *appcontext.Context, rq *dto.CreateTaskRq) (*dto.CreateTaskRs, error)
	GetByID(ctx *appcontext.Context, rq *dto.GetByIdRq) (*dto.GetByIdRs, error)
}
