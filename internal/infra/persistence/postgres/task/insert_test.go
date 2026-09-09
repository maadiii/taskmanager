package task

import (
	"context"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domain "github.com/maadiii/taskmanager/internal/domain/task"
	pkgerrors "github.com/maadiii/taskmanager/pkg/errors"
)

func TestCreateNew_Success(t *testing.T) {
	m := &MockExecer{}
	r := NewRepository(m, testTracer())

	ent := &domain.Entity{
		ID:          uuid.NewV7().String(),
		UserID:      uuid.NewV7().String(),
		Title:       "title-1",
		Description: "desc",
		Status:      domain.StatusTodo,
		Priority:    domain.PriorityMedium,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expect Exec to be called. Return a zero CommandTag and nil error.
	m.On("Exec", mock.Anything, InsertQuery, mock.Anything).Return(pgconn.CommandTag{}, nil).Run(func(args mock.Arguments) {
		// the third argument was the slice of values passed to Exec
		vals := args.Get(2).([]any)
		assert.Equal(t, ent.ID, vals[0])
		assert.Equal(t, ent.UserID, vals[1])
		assert.Equal(t, ent.Title, vals[2])
		assert.Equal(t, ent.Description, vals[3])
		assert.Equal(t, ent.Status, vals[4])
		assert.Equal(t, ent.Priority, vals[5])
		// created_at and updated_at are time.Time; ensure they are close
		assert.WithinDuration(t, ent.CreatedAt, vals[6].(time.Time), time.Second)
		assert.WithinDuration(t, ent.UpdatedAt, vals[7].(time.Time), time.Second)
	})

	err := r.CreateNew(context.Background(), ent)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestCreateNew_UniqueViolation(t *testing.T) {
	m := &MockExecer{}
	r := NewRepository(m, testTracer())

	ent := &domain.Entity{
		ID:          uuid.NewV7().String(),
		UserID:      uuid.NewV7().String(),
		Title:       "title-1",
		Description: "desc",
		Status:      domain.StatusTodo,
		Priority:    domain.PriorityMedium,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Simulate pg constraint violation
	pgErr := &pgconn.PgError{Message: "duplicate key value violates unique constraint", ConstraintName: "idx_unique_tasks_title_user_id"}
	m.On("Exec", mock.Anything, InsertQuery, mock.Anything).Return(pgconn.CommandTag{}, pgErr)

	err := r.CreateNew(context.Background(), ent)
	if assert.Error(t, err) {
		// Build an expected UniqueViolation to compare using the package Is implementation
		expected := pkgerrors.UniqueViolation(fmt.Errorf("%s", pgErr.Message), "task", "title")
		assert.True(t, pkgerrors.Is(err, expected), "expected a unique-violation error")
	}
	m.AssertExpectations(t)
}
