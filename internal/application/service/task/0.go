package task

import (
	"github.com/maadiii/taskmanager/internal/application/port"
)

type Service struct {
	repo port.TaskRepository
}

func NewService(repo port.TaskRepository) *Service {
	return &Service{repo}
}
