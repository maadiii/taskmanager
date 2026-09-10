package task

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/maadiii/taskmanager/internal/app/dto"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceUpdate_Success(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      userID,
			Permissions: []string{"update"},
		},
	}

	starter := &domaintask.Entity{
		ID:          id,
		UserID:      userID,
		Title:       "Old title",
		Description: "Old description",
		Status:      domaintask.StatusTodo,
		Priority:    domaintask.PriorityLow,
		CreatedAt:   time.Now().UTC().Add(-time.Hour),
		UpdatedAt:   time.Now().UTC().Add(-time.Minute),
	}
	baselineUpdatedAt := starter.UpdatedAt

	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})

	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).Return(starter, nil).Once()
	repo.On("UpdateByIdAndOwner", mock.Anything, mock.MatchedBy(func(entity *domaintask.Entity) bool {
		return entity != nil &&
			entity.ID == id &&
			entity.UserID == userID &&
			entity.Title == "New title" &&
			entity.Description == "New description" &&
			entity.Status == domaintask.StatusInProgress &&
			entity.Priority == domaintask.PriorityHigh &&
			entity.UpdatedAt.After(baselineUpdatedAt)
	})).Return(nil).Once()
	cache.On("Delete", mock.Anything, userID, id).Return(nil).Once()

	res, err := svc.Update(ctx, &dto.UpdateTaskRq{
		ID:          id,
		Title:       "New title",
		Description: "New description",
		Status:      string(domaintask.StatusInProgress),
		Priority:    string(domaintask.PriorityHigh),
	})
	if assert.NoError(t, err) {
		assert.Equal(t, id, res.ID)
		assert.Equal(t, "New title", res.Title)
		assert.Equal(t, "New description", res.Description)
		assert.Equal(t, domaintask.StatusInProgress.String(), res.Status)
		assert.Equal(t, domaintask.PriorityHigh.String(), res.Priority)
		assert.NotZero(t, res.UpdatedAtUnixSec)
	}

	repo.AssertExpectations(t)
}

func TestServiceUpdate_InvalidStatus(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID:      userID,
			Permissions: []string{"update"},
		},
	}

	entity := &domaintask.Entity{
		ID:          id,
		UserID:      userID,
		Title:       "Old title",
		Description: "Old description",
		Status:      domaintask.StatusTodo,
		Priority:    domaintask.PriorityLow,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})
	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).Return(entity, nil).Once()

	res, err := svc.Update(ctx, &dto.UpdateTaskRq{ID: id, Status: "NOT_A_STATUS"})
	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid status: NOT_A_STATUS")

	repo.AssertExpectations(t)
}

func TestServiceUpdate_InvalidatesCacheAfterPersistence(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: userID, Permissions: []string{"update"}},
	}
	entity := &domaintask.Entity{
		ID: id, UserID: userID, Title: "Old", Description: "Description",
		Status: domaintask.StatusTodo, Priority: domaintask.PriorityLow,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	repo := &MockTaskRepo{}
	cache := &MockTaskCache{}
	svc := newTestService(repo, &MockUoW{}, cache, testTaskMetrics{})

	repo.On("GetTaskByIdAndOwner", mock.Anything, id, userID).Return(entity, nil).Once()
	repo.On("UpdateByIdAndOwner", mock.Anything, entity).Return(nil).Once()
	cache.On("Delete", mock.Anything, userID, id).Return(nil).Once()

	_, err := svc.Update(ctx, &dto.UpdateTaskRq{ID: id, Title: "New"})

	assert.NoError(t, err)
	cache.AssertExpectations(t)
}
