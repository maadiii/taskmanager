package task

import (
	"context"
	"testing"
	"uuid"

	"github.com/maadiii/taskmanager/internal/app/dto"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceDelete_Success(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      userID,
			Permissions: []string{"delete"},
		},
	}

	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	metrics := &MockTaskMetrics{}
	metrics.On("DecTasks").Once()
	svc := newTestService(repo, &MockUoW{}, cache, metrics)
	entity := &domaintask.Entity{ID: id, UserID: userID}
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).Return(entity, nil).Once()
	repo.On("DeleteByIdAndOwner", mock.Anything, id, userID).Return(nil).Once()
	cache.On("Delete", mock.Anything, userID, id).Return(nil).Once()

	res, err := svc.Delete(ctx, &dto.DeleteTaskRq{ID: id})
	if assert.NoError(t, err) {
		assert.Equal(t, id, res.ID)
		assert.True(t, res.Deleted)
	}

	repo.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestServiceDelete_Forbidden(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      userID,
			Permissions: []string{"read"},
		},
	}

	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).
		Return(&domaintask.Entity{ID: id, UserID: userID}, nil).Once()

	res, err := svc.Delete(ctx, &dto.DeleteTaskRq{ID: id})
	assert.Nil(t, res)
	assert.EqualError(t, err, "forbidden")
	repo.AssertExpectations(t)
}

func TestServiceDelete_InvalidatesCacheAfterPersistence(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: userID, Permissions: []string{"delete"}},
	}
	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	entity := &domaintask.Entity{ID: id, UserID: userID}

	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).Return(entity, nil).Once()
	repo.On("DeleteByIdAndOwner", mock.Anything, id, userID).Return(nil).Once()
	cache.On("Delete", mock.Anything, userID, id).Return(nil).Once()

	res, err := svc.Delete(ctx, &dto.DeleteTaskRq{ID: id})

	assert.NoError(t, err)
	assert.True(t, res.Deleted)
	cache.AssertExpectations(t)
}
