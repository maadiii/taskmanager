package main

import (
	"context"
	"time"

	"github.com/maadiii/taskmanager/config"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/internal/app/service/task"
	taskCache "github.com/maadiii/taskmanager/internal/infra/cache/redis"
	"github.com/maadiii/taskmanager/internal/infra/http"
	"github.com/maadiii/taskmanager/internal/infra/persistence/postgres"
	taskPg "github.com/maadiii/taskmanager/internal/infra/persistence/postgres/task"
	"github.com/maadiii/taskmanager/pkg/observ"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.StopTimeout(10*time.Second), //nolint:mnd

		fx.Provide(context.Background),
		fx.Provide(config.NewConfig),
		provideTracing(),
		fx.Invoke(observ.InitTrace),
		providePgDb(),
		provideCache(),
		provideRepos(),
		provideDomain(),
		provideHttp(),
		fx.Invoke(http.Route),
	).Run()
}

func provideTracing() fx.Option {
	return fx.Provide(
		observ.NewExporter,
		observ.NewResource,
		observ.NewTraceProvider,
		fx.Annotate(
			func() trace.Tracer {
				return otel.Tracer("service")
			},
			fx.ResultTags(`name:"service"`),
		),
		fx.Annotate(
			func() trace.Tracer {
				return otel.Tracer("repository")
			},
			fx.ResultTags(`name:"repository"`),
		),
		fx.Annotate(
			func() trace.Tracer {
				return otel.Tracer("cache")
			},
			fx.ResultTags(`name:"cache"`),
		),
	)
}

func providePgDb() fx.Option {
	return fx.Provide(
		fx.Annotate(
			postgres.NewPostgresPool,
			fx.As(new(port.Execer)),
		),
		fx.Annotate(
			postgres.NewRepoFactory,
			fx.ParamTags(``, `name:"repository"`),
			fx.As(new(port.RepoFactory)),
		),
	)
}

func provideRepos() fx.Option {
	return fx.Provide(
		fx.Annotate(
			postgres.NewUoW,
			fx.ParamTags(``, `name:"repository"`),
		),
		fx.Annotate(
			taskPg.NewRepository,
			fx.ParamTags(``, `name:"repository"`),
			fx.As(new(port.TaskRepo)),
		),
	)
}

func provideCache() fx.Option {
	return fx.Provide(
		taskCache.NewClient,
		fx.Annotate(
			taskCache.NewTaskCache,
			fx.ParamTags(``, `name:"cache"`),
			fx.As(new(port.TaskCache)),
		),
	)
}

func provideDomain() fx.Option {
	return fx.Provide(
		fx.Annotate(
			task.NewService,
			fx.ParamTags(``, ``, ``, `name:"service"`),
			fx.As(new(port.TaskService)),
		),
	)
}

func provideHttp() fx.Option {
	return fx.Provide(
		http.NewApiGroupRouter,
	)
}
