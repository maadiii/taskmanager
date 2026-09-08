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
	svc := NewService(repo, nil)
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

	repo.On("GetTaskByIdAndOwner", mock.Anything, id, ownerID).Return(entity, nil).Once()

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
	svc := NewService(repo, nil)
	ctx := &appcontext.Context{
		Context: context.Background(),
		Identity: appcontext.Identity{
			UserID: ownerID,
		},
	}

	repo.On("GetTaskByIdAndOwner", mock.Anything, id, ownerID).Return((*domaintask.Entity)(nil), errors.New("task not found")).Once()

	res, err := svc.GetByID(ctx, &dto.GetByIdRq{ID: id})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "task not found")

	repo.AssertExpectations(t)
}

var _ port.TaskRepo = (*MockTaskRepo)(nil)
