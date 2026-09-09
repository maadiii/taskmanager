package task

import (
	"context"
	"errors"
	"testing"
	"uuid"

	domaintask "github.com/maadiii/taskmanager/internal/domain/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestList_WithoutStatus(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()
	userId := uuid.NewV7().String()
	rows := newMockRows(true, false)
	rows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			*args.Get(0).(*string) = id
			*args.Get(1).(*string) = "title"
			*args.Get(2).(*string) = "description"
			*args.Get(3).(*domaintask.Status) = domaintask.StatusTodo
		})
	rows.On("Close").Once()

	execer := new(MockExecer)
	execer.On(
		"Query",
		mock.Anything,
		"SELECT  id, title, description, status FROM tasks WHERE user_id = $1 ORDER BY id DESC LIMIT $2",
		[]any{userId, 10},
	).Return(rows, nil).Once()

	result, err := NewRepository(execer).List(context.Background(), "", userId, 10, "")

	if assert.NoError(t, err) {
		assert.Equal(t, []domaintask.Entity{{
			ID:          id,
			Title:       "title",
			Description: "description",
			Status:      domaintask.StatusTodo,
		}}, result)
	}
	execer.AssertExpectations(t)
	rows.AssertExpectations(t)
}

func TestList_WithStatusAndCursor(t *testing.T) {
	t.Parallel()

	rows := newMockRows(false)
	rows.On("Close").Once()

	execer := new(MockExecer)
	execer.On(
		"Query",
		mock.Anything,
		"SELECT  id, title, description, status FROM tasks WHERE status = $1 AND user_id = $2 AND id > $3 ORDER BY id DESC LIMIT $4",
		[]any{"DONE", "user-1", "task-10", 2},
	).Return(rows, nil).Once()

	result, err := NewRepository(execer).List(context.Background(), "DONE", "user-1", 2, "task-10")

	assert.NoError(t, err)
	assert.Empty(t, result)
	execer.AssertExpectations(t)
	rows.AssertExpectations(t)
}

func TestList_QueryError(t *testing.T) {
	t.Parallel()

	queryErr := errors.New("query failed")
	execer := new(MockExecer)
	execer.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, queryErr).Once()

	result, err := NewRepository(execer).List(context.Background(), "", "user-1", 10, "")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, queryErr)
	execer.AssertExpectations(t)
}

func TestList_ScanError(t *testing.T) {
	t.Parallel()

	scanErr := errors.New("scan failed")
	rows := newMockRows(true)
	rows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(scanErr).
		Once()
	rows.On("Close").Once()

	execer := new(MockExecer)
	execer.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(rows, nil).Once()

	result, err := NewRepository(execer).List(context.Background(), "", "user-1", 10, "")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, scanErr)
	execer.AssertExpectations(t)
	rows.AssertExpectations(t)
}

func newMockRows(next ...bool) *mockRows {
	return &mockRows{next: next}
}
