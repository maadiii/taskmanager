package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/config"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/internal/infra/persistence/postgres/task"
	"github.com/maadiii/taskmanager/pkg/errors"
	"go.uber.org/fx"
)

func NewPostgresPool(lc fx.Lifecycle, ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.PgDb.DSN)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				return errors.Wrap(err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()

			return nil
		},
	})

	return pool, nil
}

type RepoFactory struct {
	client port.Execer
}

func NewUoW(pool port.Execer) uow.UoW[port.RepoFactory] {
	return uow.NewPgx(pool, func(tx pgx.Tx) port.RepoFactory {
		return NewRepoFactory(pool)
	}).UoW()
}

func NewRepoFactory(client port.Execer) *RepoFactory {
	return &RepoFactory{
		client: client,
	}
}

func (f *RepoFactory) Tasks() port.TaskRepo {
	return task.NewRepository(f.client)
}
