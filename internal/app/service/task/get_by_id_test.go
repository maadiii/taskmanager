package task

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGetByID_Success(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()

	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID: ownerID,
		},
	}

	createdAt := time.Unix(1700000000, 0).UTC()
	entity := &domaintask.Entity{
		ID:          id,
		UserID:      ownerID,
		Title:       "Read task",
		Description: "Read task details",
		Status:      domaintask.StatusDone,
		Priority:    domaintask.PriorityHigh,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt.Add(10 * time.Minute),
	}

	cache.On("Get", mock.Anything, ownerID, id).Return((*dto.GetByIdRs)(nil), nil).Once()
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, ownerID).Return(entity, nil).Once()
	cache.On("Set", mock.Anything, ownerID, id, mock.AnythingOfType("*dto.GetByIdRs")).Return(nil).Once()

	res, err := svc.GetByID(ctx, &dto.GetByIdRq{ID: id})
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.Equal(t, entity.ID, res.ID)
		assert.Equal(t, entity.Title, res.Title)
		assert.Equal(t, entity.Description, res.Description)
		assert.Equal(t, entity.Status.String(), res.Status)
		assert.Equal(t, entity.Priority.String(), res.Priority)
		assert.Equal(t, entity.CreatedAt.Unix(), res.CreatedAtUnixSec)
	}

	repo.AssertExpectations(t)
}

func TestServiceGetByID_RepoNotFound(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()
	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID: ownerID,
		},
	}

	cache.On("Get", mock.Anything, ownerID, id).Return((*dto.GetByIdRs)(nil), nil).Once()
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, ownerID).Return((*domaintask.Entity)(nil), errors.New("task not found")).Once()

	res, err := svc.GetByID(ctx, &dto.GetByIdRq{ID: id})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "task not found")

	repo.AssertExpectations(t)
}

func TestServiceGetByID_CacheHit(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()
	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	ctx := &appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: ownerID},
	}
	cached := &dto.GetByIdRs{ID: id, Title: "Cached task"}

	cache.On("Get", mock.Anything, ownerID, id).Return(cached, nil).Once()

	res, err := svc.GetByID(ctx, &dto.GetByIdRq{ID: id})

	assert.NoError(t, err)
	assert.Equal(t, cached, res)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetTaskByIdAndOwner", mock.Anything, mock.Anything, mock.Anything)
}

func TestServiceGetByID_CacheMissStoresResult(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()
	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	ctx := &appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: ownerID},
	}
	entity := &domaintask.Entity{ID: id, UserID: ownerID, Title: "Database task"}

	cache.On("Get", mock.Anything, ownerID, id).Return((*dto.GetByIdRs)(nil), nil).Once()
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, ownerID).Return(entity, nil).Once()
	cache.On("Set", mock.Anything, ownerID, id, mock.AnythingOfType("*dto.GetByIdRs")).Return(nil).Once()

	res, err := svc.GetByID(ctx, &dto.GetByIdRq{ID: id})

	assert.NoError(t, err)
	assert.Equal(t, "Database task", res.Title)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

var _ port.TaskRepo = (*MockTaskRepo)(nil)
