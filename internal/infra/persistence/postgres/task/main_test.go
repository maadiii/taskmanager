package task

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockExecer is a testify mock implementation of port.Execer used for unit tests.
type MockExecer struct {
	mock.Mock
}

func (m *MockExecer) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	// pass the variadic arguments as a single slice into Called so expectations can inspect them
	args := m.Called(ctx, sql, arguments)
	var ct pgconn.CommandTag
	if v := args.Get(0); v != nil {
		ct = v.(pgconn.CommandTag)
	}
	return ct, args.Error(1)
}

func (m *MockExecer) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	pan := m.Called(ctx, sql, args)
	if pan.Get(0) == nil {
		return nil, pan.Error(1)
	}
	return pan.Get(0).(pgx.Rows), pan.Error(1)
}

func (m *MockExecer) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	pan := m.Called(ctx, sql, args)
	if pan.Get(0) == nil {
		return nil
	}
	return pan.Get(0).(pgx.Row)
}

func (m *MockExecer) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	pan := m.Called(ctx, txOptions)
	if pan.Get(0) == nil {
		return nil, pan.Error(1)
	}
	return pan.Get(0).(pgx.Tx), pan.Error(1)
}

// MockRow is a small pgx.Row mock used to test row.Scan behavior.
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}
