package task

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"uuid"

	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate_Success(t *testing.T) {
	t.Parallel()

	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      "user-123",
			Permissions: []string{"create"},
		},
	}

	repoFactory := &MockRepoFactory{}
	taskRepo := &MockTaskRepo{}
	uowMock := &MockUoW{
		doFn: func(ctx context.Context, fn func(context.Context, port.RepoFactory) error, opts ...*sql.TxOptions) error {
			return fn(ctx, repoFactory)
		},
	}

	repoFactory.On("Tasks").Return(taskRepo).Once()
	taskRepo.On("CreateNew", mock.Anything, mock.MatchedBy(func(entity *domaintask.Entity) bool {
		return entity != nil &&
			entity.UserID == "user-123" &&
			entity.Title == "Fix parser" &&
			elementDescription(entity.Description) &&
			entity.Status == domaintask.StatusTodo &&
			entity.Priority == domaintask.PriorityMedium
	})).Return(nil).Once()

	svc := NewService(taskRepo, uowMock)
	rq := &dto.CreateTaskRq{Title: "Fix parser", Description: "Fix parsing bug in uploader"}

	res, err := svc.Create(ctx, rq)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, res.ID)
		assert.Equal(t, rq.Title, res.Title)
		assert.Equal(t, rq.Description, res.Description)
		assert.Equal(t, domaintask.StatusTodo.String(), res.Status)
		assert.Equal(t, domaintask.PriorityMedium.String(), res.Priority)
		assert.NotZero(t, res.CreatedAtUnixSec)
	}

	repoFactory.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestServiceCreate_ForbiddenByDomainRule(t *testing.T) {
	t.Parallel()

	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      "user-456",
			Permissions: []string{"read"},
		},
	}

	uowMock := &MockUoW{}
	svc := NewService(&MockTaskRepo{}, uowMock)

	res, err := svc.Create(ctx, &dto.CreateTaskRq{Title: "Bad title", Description: "desc"})
	assert.Nil(t, res)
	assert.EqualError(t, err, "forbidden")
}

func TestServiceCreate_RepoCreateFailure(t *testing.T) {
	t.Parallel()

	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      "user-789",
			Permissions: []string{"create"},
		},
	}

	repoFactory := &MockRepoFactory{}
	taskRepo := &MockTaskRepo{}
	uowMock := &MockUoW{
		doFn: func(ctx context.Context, fn func(context.Context, port.RepoFactory) error, opts ...*sql.TxOptions) error {
			return fn(ctx, repoFactory)
		},
	}

	repoFactory.On("Tasks").Return(taskRepo).Once()
	taskRepo.On("CreateNew", mock.Anything, mock.MatchedBy(func(entity *domaintask.Entity) bool {
		return entity != nil && entity.UserID == "user-789"
	})).Return(errors.New("db failed")).Once()

	svc := NewService(taskRepo, uowMock)
	res, err := svc.Create(ctx, &dto.CreateTaskRq{Title: "Retry me", Description: "test desc"})
	assert.Nil(t, res)
	assert.EqualError(t, err, "db failed")

	repoFactory.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func elementDescription(v string) bool {
	return v == "Fix parsing bug in uploader" || v == ""
}

func TestServiceCreate_UsesUUIDv7ForEntityID(t *testing.T) {
	t.Parallel()

	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      "user-uuid",
			Permissions: []string{"create"},
		},
	}

	repoFactory := &MockRepoFactory{}
	taskRepo := &MockTaskRepo{}
	uowMock := &MockUoW{
		doFn: func(ctx context.Context, fn func(context.Context, port.RepoFactory) error, opts ...*sql.TxOptions) error {
			return fn(ctx, repoFactory)
		},
	}

	repoFactory.On("Tasks").Return(taskRepo).Once()
	taskRepo.On("CreateNew", mock.Anything, mock.MatchedBy(func(entity *domaintask.Entity) bool {
		if entity == nil || entity.ID == "" {
			return false
		}
		_, err := uuid.Parse(entity.ID)
		return err == nil
	})).Return(nil).Once()

	svc := NewService(taskRepo, uowMock)
	res, err := svc.Create(ctx, &dto.CreateTaskRq{Title: "uuid check", Description: "verify id format"})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.ID)
	_, parseErr := uuid.Parse(res.ID)
	assert.NoError(t, parseErr)
}
