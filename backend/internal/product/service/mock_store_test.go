package service

import (
	"context"

	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
)

type MockProductRepository struct {
	CreateProductFn  func(ctx context.Context, product *domain.Product) error
	ListProductsFn   func(ctx context.Context, limit int64) ([]*domain.Product, error)
	GetProductByIDFn func(ctx context.Context, productId uuid.UUID) (*domain.Product, error)
	UpdateProductFn  func(ctx context.Context, prod *domain.Product) error
	DeleteProductFn  func(ctx context.Context, productID uuid.UUID) error
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, product)
	}
	return nil
}

func (m *MockProductRepository) ListProducts(ctx context.Context, limit int64) ([]*domain.Product, error) {
	if m.ListProductsFn != nil {
		return m.ListProductsFn(ctx, limit)
	}
	return nil, nil
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, productId uuid.UUID) (*domain.Product, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, productId)
	}
	return nil, nil
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, prod *domain.Product) error {
	if m.UpdateProductFn != nil {
		return m.UpdateProductFn(ctx, prod)
	}
	return nil
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, productID)
	}
	return nil
}

type MockInventoryRepository struct {
	CreateInventoryFn func(ctx context.Context, inventory *inventoryDomain.Inventory) error
}

func (m *MockInventoryRepository) CreateInventory(ctx context.Context, inventory *inventoryDomain.Inventory) error {
	if m.CreateInventoryFn != nil {
		return m.CreateInventoryFn(ctx, inventory)
	}
	return nil
}

type MockTxStore struct {
	MockProductRepo   *MockProductRepository
	MockInventoryRepo *MockInventoryRepository
	CommitFn          func(ctx context.Context) error
	RollbackFn        func(ctx context.Context) error
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
	MockProductRepo *MockProductRepository
	BeginFn         func(ctx context.Context) (TxStore, error)
}

func (m *MockStore) Products() ProductRepository {
	return m.MockProductRepo
}

func (m *MockStore) Begin(ctx context.Context) (TxStore, error) {
	if m.BeginFn != nil {
		return m.BeginFn(ctx)
	}
	return nil, nil
}
