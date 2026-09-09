package task

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/maadiii/taskmanager/internal/app/dto"
	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceList_Success(t *testing.T) {
	t.Parallel()

	userID := uuid.NewV7().String()
	lastID := uuid.NewV7().String()
	taskID := uuid.NewV7().String()
	ctx := &appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: userID},
	}
	request := &dto.ListTaskRq{Status: "DONE", Limit: 10, LastID: lastID}
	createdAt := time.Unix(1700000000, 0).UTC()
	repo := &MockTaskRepo{}
	repo.On("List", mock.Anything, "DONE", userID, 10, lastID).Return([]domaintask.Entity{{
		ID:        taskID,
		UserID:    userID,
		Title:     "Finished task",
		Status:    domaintask.StatusDone,
		Priority:  domaintask.PriorityHigh,
		CreatedAt: createdAt,
	}}, nil).Once()

	result, err := NewService(repo, nil).List(ctx, request)

	if assert.NoError(t, err) {
		assert.Equal(t, &dto.ListTaskRs{Tasks: []dto.GetByIdRs{{
			ID:               taskID,
			Title:            "Finished task",
			Status:           "DONE",
			Priority:         "HIGH",
			CreatedAtUnixSec: createdAt.Unix(),
		}}}, result)
	}
	repo.AssertExpectations(t)
}

func TestServiceList_RepoError(t *testing.T) {
	t.Parallel()

	userID := uuid.NewV7().String()
	repoErr := errors.New("list failed")
	repo := &MockTaskRepo{}
	repo.On("List", mock.Anything, "", userID, 20, "").Return(nil, repoErr).Once()

	result, err := NewService(repo, nil).List(&appcontext.Context{
		Context:  context.Background(),
		Identity: appcontext.Identity{UserID: userID},
	}, &dto.ListTaskRq{Limit: 20})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
