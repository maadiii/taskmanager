package task

import (
	"context"
	"database/sql"

	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/dto"
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

func (m *MockTaskRepo) List(ctx context.Context, status, userID string, limit int, lastID string) ([]domaintask.Entity, error) {
	args := m.Called(ctx, status, userID, limit, lastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domaintask.Entity), args.Error(1)
}

func (m *MockTaskRepo) UpdateByIdAndOwner(ctx context.Context, entity *domaintask.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockTaskRepo) DeleteByIdAndOwner(ctx context.Context, id, ownerId string) error {
	args := m.Called(ctx, id, ownerId)
	return args.Error(0)
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

type MockTaskCache struct {
	mock.Mock
}

func newTestService(repo port.TaskRepo, uow uow.UoW[port.RepoFactory], cache port.TaskCache) *service {
	return &service{repo: repo, uow: uow, cache: cache}
}

func (m *MockTaskCache) Get(ctx context.Context, userID, taskID string) (*dto.GetByIdRs, error) {
	args := m.Called(ctx, userID, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.GetByIdRs), args.Error(1)
}

func (m *MockTaskCache) Set(ctx context.Context, userID, taskID string, task *dto.GetByIdRs) error {
	return m.Called(ctx, userID, taskID, task).Error(0)
}

func (m *MockTaskCache) Delete(ctx context.Context, userID, taskID string) error {
	return m.Called(ctx, userID, taskID).Error(0)
}

var _ uow.UoW[port.RepoFactory] = (*MockUoW)(nil)
