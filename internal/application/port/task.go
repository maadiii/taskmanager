package port

import (
	"context"

	"github.com/maadiii/taskmanager/internal/application/dto"
	"github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

type TaskRepository interface {
	CreateNewTask(ctx context.Context, entity *task.Entity) error
}

type TaskService interface {
	Create(ctx *appcontext.Context, rq *dto.CreateTaskRq) (*dto.CreateTaskRs, error)
}
