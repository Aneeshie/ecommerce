package service

import (
	"context"

	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	productDomain "github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
)

type MockOrderRepository struct {
	CreateOrderFn       func(ctx context.Context, order *domain.Order) error
	CreateOrderItemsFn  func(ctx context.Context, orderItems []*domain.OrderItem) error
	DecreaseInventoryFn func(ctx context.Context, productID uuid.UUID, quantity int) error
	GetOrdersByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetOrderByIDFn      func(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error)
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	if m.CreateOrderFn != nil {
		return m.CreateOrderFn(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) CreateOrderItems(ctx context.Context, orderItems []*domain.OrderItem) error {
	if m.CreateOrderItemsFn != nil {
		return m.CreateOrderItemsFn(ctx, orderItems)
	}
	return nil
}

func (m *MockOrderRepository) DecreaseInventory(ctx context.Context, productID uuid.UUID, quantity int) error {
	if m.DecreaseInventoryFn != nil {
		return m.DecreaseInventoryFn(ctx, productID, quantity)
	}
	return nil
}

func (m *MockOrderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	if m.GetOrdersByUserIDFn != nil {
		return m.GetOrdersByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockOrderRepository) GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error) {
	if m.GetOrderByIDFn != nil {
		return m.GetOrderByIDFn(ctx, userID, orderID)
	}
	return nil, nil
}

type MockProductRepository struct {
	GetProductByIDFn func(ctx context.Context, productId uuid.UUID) (*productDomain.Product, error)
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, productId uuid.UUID) (*productDomain.Product, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, productId)
	}
	return nil, nil
}

type MockInventoryRepository struct {
	GetInventoryByProductIDFn func(ctx context.Context, productID uuid.UUID) (*inventoryDomain.Inventory, error)
}

func (m *MockInventoryRepository) GetInventoryByProductID(ctx context.Context, productID uuid.UUID) (*inventoryDomain.Inventory, error) {
	if m.GetInventoryByProductIDFn != nil {
		return m.GetInventoryByProductIDFn(ctx, productID)
	}
	return nil, nil
}

type MockTxStore struct {
	MockOrderRepo     *MockOrderRepository
	MockProductRepo   *MockProductRepository
	MockInventoryRepo *MockInventoryRepository
	CommitFn          func(ctx context.Context) error
	RollbackFn        func(ctx context.Context) error
}

func (m *MockTxStore) Orders() OrderRepository {
	return m.MockOrderRepo
}

func (m *MockTxStore) Products() ProductRepository {
	return m.MockProductRepo
}

func (m *MockTxStore) Inventory() InventoryRepository {
	return m.MockInventoryRepo
}

func (m *MockTxStore) Commit(ctx context.Context) error {
	if m.CommitFn != nil {
		return m.CommitFn(ctx)
	}
	return nil
}

func (m *MockTxStore) Rollback(ctx context.Context) error {
	if m.RollbackFn != nil {
		return m.RollbackFn(ctx)
	}
	return nil
}

type MockStore struct {
	MockOrderRepo     *MockOrderRepository
	MockProductRepo   *MockProductRepository
	MockInventoryRepo *MockInventoryRepository
	BeginFn           func(ctx context.Context) (TxStore, error)
}

func (m *MockStore) Orders() OrderRepository {
	return m.MockOrderRepo
}

func (m *MockStore) Products() ProductRepository {
	return m.MockProductRepo
}

func (m *MockStore) Inventory() InventoryRepository {
	return m.MockInventoryRepo
}

func (m *MockStore) Begin(ctx context.Context) (TxStore, error) {
	if m.BeginFn != nil {
		return m.BeginFn(ctx)
	}
	return nil, nil
}
