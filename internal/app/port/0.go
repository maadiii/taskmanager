package port

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maadiii/taskmanager/internal/app/dto"
)

type TaskCache interface {
	Get(ctx context.Context, userID, taskID string) (*dto.GetByIdRs, error)
	Set(ctx context.Context, userID, taskID string, task *dto.GetByIdRs) error
	Delete(ctx context.Context, userID, taskID string) error
}

type TaskMetrics interface {
	IncTasks()
	DecTasks()
}

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type RepoFactory interface {
	Tasks() TaskRepo
}
