package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/common/testdb"
	"github.com/Aneeshie/ecommerce/internal/inventory"
	"github.com/Aneeshie/ecommerce/internal/inventory/domain"
	productDomain "github.com/Aneeshie/ecommerce/internal/product/domain"
	productRepo "github.com/Aneeshie/ecommerce/internal/product/repository"
	"github.com/google/uuid"
)

func TestCreateInventory(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)
	prodRepo := productRepo.NewRepository(tx)

	now := time.Now()
	amount, _ := money.New(100)
	prod := &productDomain.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "A test product",
		Price:       amount,
		Status:      productDomain.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := prodRepo.CreateProduct(context.Background(), prod)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	inv := &domain.Inventory{
		ProductID: prod.ID,
		Quantity:  10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = repo.CreateInventory(context.Background(), inv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Fetch to verify
	fetched, err := repo.GetInventoryByProductID(context.Background(), inv.ProductID)
	if err != nil {
		t.Fatalf("failed to fetch created inventory: %v", err)
	}

	if fetched.Quantity != inv.Quantity {
		t.Errorf("expected quantity %d, got %d", inv.Quantity, fetched.Quantity)
	}
}

func TestGetInventoryByProductID_NotFound(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	_, err := repo.GetInventoryByProductID(context.Background(), uuid.New())
	if err != inventory.ErrNoProductFound {
		t.Fatalf("expected ErrNoProductFound, got %v", err)
	}
}
