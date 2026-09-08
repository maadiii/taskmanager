package task

import (
	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/port"
)

type service struct {
	repo port.TaskRepo
	uow  uow.UoW[port.RepoFactory]
}

func NewService(
	repo port.TaskRepo,
	uow uow.UoW[port.RepoFactory],
) *service {
	return &service{repo, uow}
}
