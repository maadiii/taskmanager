package task

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	domain "github.com/maadiii/taskmanager/internal/domain/task"
	pkgerrors "github.com/maadiii/taskmanager/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTaskByIdAndOwner_Success(t *testing.T) {
	t.Parallel()

	m := &MockExecer{}
	r := NewRepository(m)

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedAt := createdAt.Add(30 * time.Minute)

	row := &MockRow{}
	row.On(
		"Scan",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil).Run(func(args mock.Arguments) {
		dests := []any{
			args.Get(0),
			args.Get(1),
			args.Get(2),
			args.Get(3),
			args.Get(4),
			args.Get(5),
			args.Get(6),
		}

		*dests[0].(*string) = id
		*dests[1].(*string) = "task-title"
		*dests[2].(*string) = "task-description"
		*dests[3].(*domain.Status) = domain.StatusInProgress
		*dests[4].(*domain.Priority) = domain.PriorityHigh
		*dests[5].(*time.Time) = createdAt
		*dests[6].(*time.Time) = updatedAt
	})

	m.On("QueryRow", mock.Anything, GetByIdAndOwnerQuery, []any{id, ownerID}).Return(row)

	entity, err := r.GetTaskByIdAndOwner(context.Background(), id, ownerID)
	if assert.NoError(t, err) {
		assert.Equal(t, id, entity.ID)
		assert.Equal(t, "task-title", entity.Title)
		assert.Equal(t, "task-description", entity.Description)
		assert.Equal(t, domain.StatusInProgress, entity.Status)
		assert.Equal(t, domain.PriorityHigh, entity.Priority)
		assert.WithinDuration(t, createdAt, entity.CreatedAt, time.Second)
		assert.WithinDuration(t, updatedAt, entity.UpdatedAt, time.Second)
	}

	m.AssertExpectations(t)
	row.AssertExpectations(t)
}

func TestGetTaskByIdAndOwner_NotFound(t *testing.T) {
	t.Parallel()

	m := &MockExecer{}
	r := NewRepository(m)

	id := uuid.NewV7().String()
	ownerID := uuid.NewV7().String()

	row := &MockRow{}
	row.On(
		"Scan",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(pgx.ErrNoRows)

	m.On("QueryRow", mock.Anything, GetByIdAndOwnerQuery, []any{id, ownerID}).Return(row)

	_, err := r.GetTaskByIdAndOwner(context.Background(), id, ownerID)
	if assert.Error(t, err) {
		expected := pkgerrors.NotFound(pgx.ErrNoRows, "task_id")
		assert.True(t, pkgerrors.Is(err, expected), "expected a not-found task error")
	}

	m.AssertExpectations(t)
	row.AssertExpectations(t)
}
