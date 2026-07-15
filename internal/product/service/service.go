package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/Aneeshie/ecommerce/internal/store"
	"github.com/google/uuid"
)

var (
	ErrEmptyProductName        = errors.New("Product name cannot be empty")
	ErrEmptyProductDescription = errors.New("Product description cannot be empty")
)

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	amount, err := money.New(req.Price)
	if err != nil {
		return nil, money.ErrNegativeAmount
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrEmptyProductName
	}

	if strings.TrimSpace(req.Description) == "" {
		return nil, ErrEmptyProductDescription
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

	inventory := &inventoryDomain.Inventory {
		ProductID: product.ID,
		Quantity: 0,
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

func (s *Service) UpdateProduct(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
	existing, err := s.store.Products().GetProductByID(ctx, productID)
	if err != nil {
		return err

	}
	if strings.TrimSpace(product.Name) == "" {
		return ErrEmptyProductName
	}

	if strings.TrimSpace(product.Description) == "" {
		return ErrEmptyProductDescription
	}

	amount, err := money.New(product.Price)
	if err != nil {
		return err
	}

	existing.Name = product.Name
	existing.Description = product.Description
	existing.Price = amount
	existing.UpdatedAt = time.Now()

	return s.store.Products().UpdateProduct(ctx, existing)
}

func (s *Service) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.store.Products().DeleteProduct(ctx, productID)
}
