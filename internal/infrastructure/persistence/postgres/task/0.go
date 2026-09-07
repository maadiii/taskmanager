package task

import (
	"github.com/maadiii/taskmanager/internal/application/port"
)

type Repository struct {
	client port.Execer
}

func NewRepository(client port.Execer) *Repository {
	return &Repository{client}
}
