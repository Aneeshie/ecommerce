package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/order"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	productDomain "github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
)

func TestCreateOrder(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	amount, _ := money.New(1000)

	tests := []struct {
		Name          string
		Req           *dto.CreateOrderRequest
		SetupMock     func(m *MockStore)
		ExpectedError error
	}{
		{
			Name: "Successful Create",
			Req: &dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: productID, Quantity: 2},
				},
			},
			SetupMock: func(m *MockStore) {
				m.MockProductRepo = &MockProductRepository{
					GetProductByIDFn: func(ctx context.Context, id uuid.UUID) (*productDomain.Product, error) {
						return &productDomain.Product{ID: id, Price: amount}, nil
					},
				}
				m.MockInventoryRepo = &MockInventoryRepository{
					GetInventoryByProductIDFn: func(ctx context.Context, id uuid.UUID) (*inventoryDomain.Inventory, error) {
						return &inventoryDomain.Inventory{ProductID: id, Quantity: 10}, nil
					},
				}
				m.BeginFn = func(ctx context.Context) (TxStore, error) {
					return &MockTxStore{
						MockOrderRepo: &MockOrderRepository{},
					}, nil
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Empty Order",
			Req: &dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{},
			},
			SetupMock:     func(m *MockStore) {},
			ExpectedError: order.ErrEmptyOrder,
		},
		{
			Name: "Insufficient Inventory",
			Req: &dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: productID, Quantity: 5},
				},
			},
			SetupMock: func(m *MockStore) {
				m.MockProductRepo = &MockProductRepository{
					GetProductByIDFn: func(ctx context.Context, id uuid.UUID) (*productDomain.Product, error) {
						return &productDomain.Product{ID: id, Price: amount}, nil
					},
				}
				m.MockInventoryRepo = &MockInventoryRepository{
					GetInventoryByProductIDFn: func(ctx context.Context, id uuid.UUID) (*inventoryDomain.Inventory, error) {
						return &inventoryDomain.Inventory{ProductID: id, Quantity: 2}, nil
					},
				}
			},
			ExpectedError: order.ErrInsufficientInventory,
		},
		{
			Name: "Duplicate Product",
			Req: &dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: productID, Quantity: 2},
					{ProductID: productID, Quantity: 1},
				},
			},
			SetupMock:     func(m *MockStore) {},
			ExpectedError: order.ErrDuplicateProduct,
		},
		{
			Name: "Transaction Begin Error",
			Req: &dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: productID, Quantity: 2},
				},
			},
			SetupMock: func(m *MockStore) {
				m.MockProductRepo = &MockProductRepository{
					GetProductByIDFn: func(ctx context.Context, id uuid.UUID) (*productDomain.Product, error) {
						return &productDomain.Product{ID: id, Price: amount}, nil
					},
				}
				m.MockInventoryRepo = &MockInventoryRepository{
					GetInventoryByProductIDFn: func(ctx context.Context, id uuid.UUID) (*inventoryDomain.Inventory, error) {
						return &inventoryDomain.Inventory{ProductID: id, Quantity: 10}, nil
					},
				}
				m.BeginFn = func(ctx context.Context) (TxStore, error) {
					return nil, errors.New("begin error")
				}
			},
			ExpectedError: errors.New("begin error"),
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

			err := service.CreateOrder(context.Background(), userID, tt.Req)

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

func TestGetOrdersByUserID(t *testing.T) {
	userID := uuid.New()
	amount, _ := money.New(100)
	ords := []*domain.Order{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Status:     domain.StatusPending,
			TotalPrice: amount,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
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
				m.MockOrderRepo = &MockOrderRepository{
					GetOrdersByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]*domain.Order, error) {
						return ords, nil
					},
				}
			},
			ExpectedError: nil,
		},
		{
			Name: "Repository Error",
			SetupMock: func(m *MockStore) {
				m.MockOrderRepo = &MockOrderRepository{
					GetOrdersByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]*domain.Order, error) {
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

			_, err := service.GetOrdersByUserID(context.Background(), userID)

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
