package service

import (
	"context"
	"strings"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/product"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/Aneeshie/ecommerce/internal/store"
	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	ListProducts(ctx context.Context, limit int64) ([]*domain.Product, error)
	GetProductByID(ctx context.Context, productId uuid.UUID) (*domain.Product, error)
	UpdateProduct(ctx context.Context, prod *domain.Product) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}

type InventoryRepository interface {
	CreateInventory(ctx context.Context, inventory *inventoryDomain.Inventory) error
}

type TxStore interface {
	Products() ProductRepository
	Inventory() InventoryRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Store interface {
	Products() ProductRepository
	Begin(ctx context.Context) (TxStore, error)
}

type storeWrapper struct {
	*store.Store
}

func (w *storeWrapper) Products() ProductRepository {
	return w.Store.Products()
}

func (w *storeWrapper) Begin(ctx context.Context) (TxStore, error) {
	tx, err := w.Store.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx}, nil
}

type txWrapper struct {
	*store.TxStore
}

func (t *txWrapper) Products() ProductRepository {
	return t.TxStore.Products()
}

func (t *txWrapper) Inventory() InventoryRepository {
	return t.TxStore.Inventory()
}

type Service struct {
	store Store
}

func NewService(s *store.Store) *Service {
	return &Service{
		store: &storeWrapper{s},
	}
}

func (s *Service) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	amount, err := money.New(req.Price)
	if err != nil {
		return nil, money.ErrNegativeAmount
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, product.ErrEmptyProductName
	}

	if strings.TrimSpace(req.Description) == "" {
		return nil, product.ErrEmptyProductDescription
	}

	product := domain.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       amount,
		Status:      domain.ProductStatusActive,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.store.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = tx.Products().CreateProduct(ctx, &product)
	if err != nil {
		return &dto.CreateProductResponse{}, err
	}

	inventory := &inventoryDomain.Inventory{
		ProductID: product.ID,
		Quantity:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = tx.Inventory().CreateInventory(ctx, inventory)
	if err != nil {
		return &dto.CreateProductResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return &dto.CreateProductResponse{}, err
	}

	return &dto.CreateProductResponse{
		ID: product.ID,
	}, nil
}

func (s *Service) ListProducts(ctx context.Context, limit int64) ([]*dto.ProductResponse, error) {
	products, err := s.store.Products().ListProducts(ctx, limit)
	if err != nil {
		return nil, err
	}

	var response []*dto.ProductResponse

	for _, product := range products {
		response = append(response, &dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price.Amount(),
			Status:      string(product.Status),
		})
	}

	return response, nil
}

func (s *Service) GetProductById(ctx context.Context, productId uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.store.Products().GetProductByID(ctx, productId)
	if err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.Amount(),
		Status:      string(product.Status),
	}, nil
}

func (s *Service) UpdateProduct(ctx context.Context, productID uuid.UUID, prod *dto.UpdateProductRequest) error {
	existing, err := s.store.Products().GetProductByID(ctx, productID)
	if err != nil {
		return err

	}
	if strings.TrimSpace(prod.Name) == "" {
		return product.ErrEmptyProductName
	}

	if strings.TrimSpace(prod.Description) == "" {
		return product.ErrEmptyProductDescription
	}

	amount, err := money.New(prod.Price)
	if err != nil {
		return err
	}

	existing.Name = prod.Name
	existing.Description = prod.Description
	existing.Price = amount
	existing.UpdatedAt = time.Now()

	return s.store.Products().UpdateProduct(ctx, existing)
}

func (s *Service) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.store.Products().DeleteProduct(ctx, productID)
}
