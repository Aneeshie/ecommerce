package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/inventory"
	"github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateInventory(t *testing.T) {
	tests := []struct {
		Name          string
		Inventory     *domain.Inventory
		SetupMock     func(m *MockQueryExecutor)
		ExpectedError error
	}{
		{
			Name: "Successful Create",
			Inventory: &domain.Inventory{
				ProductID: uuid.New(),
				Quantity:  10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			SetupMock: func(m *MockQueryExecutor) {
				m.ExecFn = func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
					return pgconn.NewCommandTag("INSERT 0 1"), nil
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Database Error",
			Inventory: &domain.Inventory{
				ProductID: uuid.New(),
				Quantity:  10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			SetupMock: func(m *MockQueryExecutor) {
				m.ExecFn = func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
					return pgconn.CommandTag{}, errors.New("db error")
				}
			},
			ExpectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB := &MockQueryExecutor{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockDB)
			}
			repo := NewRepository(mockDB)

			err := repo.CreateInventory(context.Background(), tt.Inventory)

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

func TestGetInventoryByProductID(t *testing.T) {
	productID := uuid.New()
	now := time.Now()

	tests := []struct {
		Name          string
		ProductID     uuid.UUID
		SetupMock     func(m *MockQueryExecutor)
		ExpectedError error
	}{
		{
			Name:      "Successful Get",
			ProductID: productID,
			SetupMock: func(m *MockQueryExecutor) {
				m.QueryRowFn = func(ctx context.Context, sql string, args ...any) pgx.Row {
					return &MockRow{
						ScanFn: func(dest ...any) error {
							*dest[0].(*uuid.UUID) = productID
							*dest[1].(*int) = 10
							*dest[2].(*time.Time) = now
							*dest[3].(*time.Time) = now
							return nil
						},
					}
				}
			},
			ExpectedError: nil,
		},
		{
			Name:      "Not Found",
			ProductID: productID,
			SetupMock: func(m *MockQueryExecutor) {
				m.QueryRowFn = func(ctx context.Context, sql string, args ...any) pgx.Row {
					return &MockRow{
						ScanFn: func(dest ...any) error {
							return pgx.ErrNoRows
						},
					}
				}
			},
			ExpectedError: inventory.ErrNoProductFound,
		},
		{
			Name:      "Database Error",
			ProductID: productID,
			SetupMock: func(m *MockQueryExecutor) {
				m.QueryRowFn = func(ctx context.Context, sql string, args ...any) pgx.Row {
					return &MockRow{
						ScanFn: func(dest ...any) error {
							return errors.New("query error")
						},
					}
				}
			},
			ExpectedError: errors.New("query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB := &MockQueryExecutor{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockDB)
			}
			repo := NewRepository(mockDB)

			_, err := repo.GetInventoryByProductID(context.Background(), tt.ProductID)

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
