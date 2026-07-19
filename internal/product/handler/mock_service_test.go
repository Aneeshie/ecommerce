package handler

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/google/uuid"
)

type MockProductService struct {
	CreateProductFn  func(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error)
	ListProductsFn   func(ctx context.Context, limit int64) ([]*dto.ProductResponse, error)
	GetProductByIdFn func(ctx context.Context, productId uuid.UUID) (*dto.ProductResponse, error)
	UpdateProductFn  func(ctx context.Context, productID uuid.UUID, prod *dto.UpdateProductRequest) error
	DeleteProductFn  func(ctx context.Context, productID uuid.UUID) error
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, req)
	}
	return nil, nil
}

func (m *MockProductService) ListProducts(ctx context.Context, limit int64) ([]*dto.ProductResponse, error) {
	if m.ListProductsFn != nil {
		return m.ListProductsFn(ctx, limit)
	}
	return nil, nil
}

func (m *MockProductService) GetProductById(ctx context.Context, productId uuid.UUID) (*dto.ProductResponse, error) {
	if m.GetProductByIdFn != nil {
		return m.GetProductByIdFn(ctx, productId)
	}
	return nil, nil
}

func (m *MockProductService) UpdateProduct(ctx context.Context, productID uuid.UUID, prod *dto.UpdateProductRequest) error {
	if m.UpdateProductFn != nil {
		return m.UpdateProductFn(ctx, productID, prod)
	}
	return nil
}

func (m *MockProductService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, productID)
	}
	return nil
}
