package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/product"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/google/uuid"
)

func TestCreateProduct(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:        "Test",
		Description: "Desc",
		Price:       1000,
	}

	tests := []struct {
		Name          string
		Req           *dto.CreateProductRequest
		SetupMock     func(m *MockStore)
		ExpectedError error
	}{
		{
			Name: "Successful Create",
			Req:  req,
			SetupMock: func(m *MockStore) {
				tx := &MockTxStore{
					MockProductRepo:   &MockProductRepository{},
					MockInventoryRepo: &MockInventoryRepository{},
				}
				m.BeginFn = func(ctx context.Context) (TxStore, error) {
					return tx, nil
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Empty Name",
			Req: &dto.CreateProductRequest{
				Name:        "",
				Description: "Desc",
				Price:       1000,
			},
			SetupMock:     func(m *MockStore) {},
			ExpectedError: product.ErrEmptyProductName,
		},
		{
			Name: "Empty Description",
			Req: &dto.CreateProductRequest{
				Name:        "Test",
				Description: "",
				Price:       1000,
			},
			SetupMock:     func(m *MockStore) {},
			ExpectedError: product.ErrEmptyProductDescription,
		},
		{
			Name: "Negative Price",
			Req: &dto.CreateProductRequest{
				Name:        "Test",
				Description: "Desc",
				Price:       -100,
			},
			SetupMock:     func(m *MockStore) {},
			ExpectedError: money.ErrNegativeAmount,
		},
		{
			Name: "Tx Begin Error",
			Req:  req,
			SetupMock: func(m *MockStore) {
				m.BeginFn = func(ctx context.Context) (TxStore, error) {
					return nil, errors.New("begin error")
				}
			},
			ExpectedError: errors.New("begin error"),
		},
		{
			Name: "Create Inventory Error",
			Req:  req,
			SetupMock: func(m *MockStore) {
				tx := &MockTxStore{
					MockProductRepo: &MockProductRepository{},
					MockInventoryRepo: &MockInventoryRepository{
						CreateInventoryFn: func(ctx context.Context, inventory *inventoryDomain.Inventory) error {
							return errors.New("inventory error")
						},
					},
				}
				m.BeginFn = func(ctx context.Context) (TxStore, error) {
					return tx, nil
				}
			},
			ExpectedError: errors.New("inventory error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockStore := &MockStore{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockStore)
			}
			service := &Service{
				store: mockStore,
			}

			_, err := service.CreateProduct(context.Background(), tt.Req)

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

func TestListProducts(t *testing.T) {
	amount, _ := money.New(100)
	prods := []*domain.Product{
		{
			ID:          uuid.New(),
			Name:        "Test",
			Description: "Desc",
			Price:       amount,
			Status:      domain.ProductStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tests := []struct {
		Name          string
		SetupMock     func(m *MockStore)
		ExpectedError error
	}{
		{
			Name: "Successful List",
			SetupMock: func(m *MockStore) {
				m.MockProductRepo = &MockProductRepository{
					ListProductsFn: func(ctx context.Context, limit int64) ([]*domain.Product, error) {
						return prods, nil
					},
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Repository Error",
			SetupMock: func(m *MockStore) {
				m.MockProductRepo = &MockProductRepository{
					ListProductsFn: func(ctx context.Context, limit int64) ([]*domain.Product, error) {
						return nil, errors.New("repo error")
					},
				}
			},
			ExpectedError: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockStore := &MockStore{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockStore)
			}
			service := &Service{
				store: mockStore,
			}

			_, err := service.ListProducts(context.Background(), 10)

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
