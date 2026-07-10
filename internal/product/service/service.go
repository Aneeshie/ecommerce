package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/Aneeshie/ecommerce/internal/product/repository"
	"github.com/google/uuid"
)

var (
	ErrEmptyProductName = errors.New("Product name cannot be empty")
	ErrEmptyProductDescription =  errors.New("Product description cannot be empty")
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service{
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	amount, err  := money.New(req.Price)
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
		ID: uuid.New(),
		Name: req.Name,
		Description: req.Description,
		Price: amount,
		Status: domain.ProductStatusActive,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateProduct(ctx, &product)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProductResponse{
		ID: product.ID,
	}, nil
}

func (s *Service) ListProducts (ctx context.Context, limit int64) ([]*dto.ProductResponse, error){
	products, err := s.repo.ListProducts(ctx, limit)
	if err != nil {
		return nil, err
	}

	var response []*dto.ProductResponse

	for _, product := range products {
		response = append(response, &dto.ProductResponse{
			ID: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: product.Price.Amount(),
			Status: string(product.Status),
		})
	}

	return response, nil
}
