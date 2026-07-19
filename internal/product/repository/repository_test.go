package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/product"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateProduct(t *testing.T) {
	amount, _ := money.New(1000)
	prod := &domain.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "A nice product",
		Price:       amount,
		Status:      domain.ProductStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		Name          string
		Product       *domain.Product
		SetupMock     func(m *MockQueryExecutor)
		ExpectedError error
	}{
		{
			Name:    "Successful Create",
			Product: prod,
			SetupMock: func(m *MockQueryExecutor) {
				m.ExecFn = func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
					return pgconn.NewCommandTag("INSERT 0 1"), nil
				}
			},
			ExpectedError: nil,
		},
		{
			Name:    "Database Error",
			Product: prod,
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

			err := repo.CreateProduct(context.Background(), tt.Product)

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

func TestGetProductByID(t *testing.T) {
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
							*dest[1].(*string) = "Test Product"
							*dest[2].(*string) = "A nice product"
							*dest[3].(*int64) = 1000
							*dest[4].(*domain.Status) = domain.ProductStatusActive
							*dest[5].(*time.Time) = now
							*dest[6].(*time.Time) = now
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
			ExpectedError: product.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB := &MockQueryExecutor{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockDB)
			}
			repo := NewRepository(mockDB)

			_, err := repo.GetProductByID(context.Background(), tt.ProductID)

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
