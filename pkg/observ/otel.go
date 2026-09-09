package observ

import (
	"context"

	"github.com/maadiii/taskmanager/config"
	"github.com/maadiii/taskmanager/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.uber.org/fx"
)

func NewExporter(ctx context.Context, cfg *config.Config) (*otlptrace.Exporter, error) {
	exporter, err := otlptracehttp.New(
		ctx, otlptracehttp.WithEndpoint(cfg.Server.OtlpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return exporter, nil
}

func NewResource(ctx context.Context, cfg *config.Config) (*resource.Resource, error) {
	res, err := resource.New(
		ctx, resource.WithAttributes(semconv.ServiceName(cfg.Server.Name)),
	)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return res, nil
}

func NewTraceProvider(exporter *otlptrace.Exporter, resource *resource.Resource) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
	)
}

func InitTrace(lc fx.Lifecycle, tp *trace.TracerProvider) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			))

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := tp.Shutdown(ctx); err != nil {
				return errors.Wrap(err)
			}

			return nil
		},
	})
}
