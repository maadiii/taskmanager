package task

import (
	"context"
	"database/sql"

	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/port"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateNew(ctx context.Context, entity *domaintask.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockTaskRepo) GetTaskByIdAndOwner(ctx context.Context, id, ownerId string) (*domaintask.Entity, error) {
	args := m.Called(ctx, id, ownerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domaintask.Entity), args.Error(1)
}

type MockRepoFactory struct {
	mock.Mock
}

func (m *MockRepoFactory) Tasks() port.TaskRepo {
	args := m.Called()
	return args.Get(0).(port.TaskRepo)
}

type MockUoW struct {
	doFn func(ctx context.Context, fn func(ctx context.Context, repo port.RepoFactory) error, opts ...*sql.TxOptions) error
}

func (m *MockUoW) Do(ctx context.Context, fn func(ctx context.Context, repo port.RepoFactory) error, opts ...*sql.TxOptions) error {
	if m.doFn != nil {
		return m.doFn(ctx, fn, opts...)
	}
	return nil
}

func (m *MockUoW) Begin(ctx context.Context, opts ...*sql.TxOptions) (context.Context, port.RepoFactory, error) {
	return ctx, nil, nil
}

func (m *MockUoW) SavePoint(ctx context.Context, name string) error {
	return nil
}

func (m *MockUoW) Commit(ctx context.Context) error {
	return nil
}

func (m *MockUoW) Rollback(ctx context.Context) error {
	return nil
}

var _ uow.UoW[port.RepoFactory] = (*MockUoW)(nil)
