package main

import (
	"context"
	"time"

	"github.com/maadiii/taskmanager/config"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/internal/app/service/task"
	"github.com/maadiii/taskmanager/internal/infra/http"
	"github.com/maadiii/taskmanager/internal/infra/persistence/postgres"
	taskPg "github.com/maadiii/taskmanager/internal/infra/persistence/postgres/task"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Gracefully shutting down
		fx.StopTimeout(10*time.Second), //nolint:mnd

		fx.Provide(context.Background),
		provieConfig(),
		providePgDb(),
		provideRepos(),
		provideDomain(),
		provideHttp(),
		fx.Invoke(http.Route),
	).Run()
}

func provieConfig() fx.Option {
	return fx.Provide(
		config.NewConfig,
	)
}

func providePgDb() fx.Option {
	return fx.Provide(
		// Inject postgres pgx pool as port.Execer
		fx.Annotate(
			postgres.NewPostgresPool,
			fx.As(new(port.Execer)),
		),

		// Inject *postgres.RepoFactory as port.RepoFactory
		fx.Annotate(
			postgres.NewRepoFactory,
			fx.As(new(port.RepoFactory)),
		),
	)
}

func provideRepos() fx.Option {
	return fx.Provide(
		postgres.NewUoW,

		fx.Annotate(
			taskPg.NewRepository,
			fx.As(new(port.TaskRepo)),
		),
	)
}

func provideDomain() fx.Option {
	return fx.Provide(
		// Inject task.Service as port.TaskService
		fx.Annotate(
			task.NewService,
			fx.As(new(port.TaskService)),
		),
	)
}

func provideHttp() fx.Option {
	return fx.Provide(
		http.NewApiGroupRouter,
	)
}
