package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/order"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateOrder(t *testing.T) {
	amount, _ := money.New(1000)
	ord := &domain.Order{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Status:     domain.StatusPending,
		TotalPrice: amount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		Name          string
		Order         *domain.Order
		SetupMock     func(m *MockQueryExecutor)
		ExpectedError error
	}{
		{
			Name:  "Successful Create",
			Order: ord,
			SetupMock: func(m *MockQueryExecutor) {
				m.ExecFn = func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
					return pgconn.NewCommandTag("INSERT 0 1"), nil
				}
			},
			ExpectedError: nil,
		},
		{
			Name:  "Database Error",
			Order: ord,
			SetupMock: func(m *MockQueryExecutor) {
				m.ExecFn = func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
					return pgconn.CommandTag{}, errors.New("db error")
				}
			},
			ExpectedError: errors.New("create order: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB := &MockQueryExecutor{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockDB)
			}
			repo := NewRepository(mockDB)

			err := repo.CreateOrder(context.Background(), tt.Order)

			if tt.ExpectedError != nil {
				if err == nil || err.Error() != tt.ExpectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.ExpectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected nil error got %v", err)
			}
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		Name          string
		SetupMock     func(m *MockQueryExecutor)
		ExpectedError error
	}{
		{
			Name: "Successful Get",
			SetupMock: func(m *MockQueryExecutor) {
				m.QueryRowFn = func(ctx context.Context, sql string, args ...any) pgx.Row {
					return &MockRow{
						ScanFn: func(dest ...any) error {
							*dest[0].(*uuid.UUID) = orderID
							*dest[1].(*uuid.UUID) = userID
							*dest[2].(*domain.OrderStatus) = domain.StatusPending
							*dest[3].(*int64) = 1000
							*dest[4].(*time.Time) = now
							*dest[5].(*time.Time) = now
							return nil
						},
					}
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Not Found",
			SetupMock: func(m *MockQueryExecutor) {
				m.QueryRowFn = func(ctx context.Context, sql string, args ...any) pgx.Row {
					return &MockRow{
						ScanFn: func(dest ...any) error {
							return pgx.ErrNoRows
						},
					}
				}
			},
			ExpectedError: order.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB := &MockQueryExecutor{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockDB)
			}
			repo := NewRepository(mockDB)

			_, err := repo.GetOrderByID(context.Background(), userID, orderID)

			if tt.ExpectedError != nil {
				if err == nil || err.Error() != tt.ExpectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.ExpectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected nil error got %v", err)
			}
		})
	}
}
