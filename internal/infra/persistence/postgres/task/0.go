package task

import "github.com/maadiii/taskmanager/internal/app/port"

type repo struct {
	client port.Execer
}

func NewRepository(client port.Execer) *repo {
	return &repo{client}
}
