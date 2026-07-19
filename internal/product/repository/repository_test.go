package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/common/testdb"
	"github.com/Aneeshie/ecommerce/internal/product"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
)

func TestCreateAndGetProduct(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	amount, _ := money.New(1500)
	now := time.Now()
	prod := &domain.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "A description",
		Price:       amount,
		Status:      domain.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repo.CreateProduct(context.Background(), prod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fetched, err := repo.GetProductByID(context.Background(), prod.ID)
	if err != nil {
		t.Fatalf("failed to fetch product: %v", err)
	}

	if fetched.Name != prod.Name {
		t.Errorf("expected name %s, got %s", prod.Name, fetched.Name)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	_, err := repo.GetProductByID(context.Background(), uuid.New())
	if err != product.ErrProductNotFound {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
