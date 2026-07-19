package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockQueryExecutor struct {
	ExecFn     func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryFn    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRowFn func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (m *MockQueryExecutor) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.ExecFn != nil {
		return m.ExecFn(ctx, sql, arguments...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockQueryExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFn != nil {
		return m.QueryFn(ctx, sql, args...)
	}
	return nil, nil
}

func (m *MockQueryExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFn != nil {
		return m.QueryRowFn(ctx, sql, args...)
	}
	return &MockRow{}
}

type MockRow struct {
	ScanFn func(dest ...any) error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.ScanFn != nil {
		return m.ScanFn(dest...)
	}
	return nil
}
