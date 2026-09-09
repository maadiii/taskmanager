package task

import (
	"github.com/maadiii/taskmanager/internal/app/port"
	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	client port.Execer
	tracer trace.Tracer
}

func NewRepository(client port.Execer, tracer trace.Tracer) *repo {
	return &repo{client, tracer}
}
