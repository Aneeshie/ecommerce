package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/common/testdb"
	identityDomain "github.com/Aneeshie/ecommerce/internal/identity/domain"
	identityRepo "github.com/Aneeshie/ecommerce/internal/identity/repository"
	"github.com/Aneeshie/ecommerce/internal/order"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/google/uuid"
)

func TestCreateAndGetOrder(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)
	userRepo := identityRepo.NewRepository(tx)

	amount, _ := money.New(1000)
	now := time.Now()
	userID := uuid.New()

	user := identityDomain.User{
		ID:           userID,
		Email:        "order_test@example.com",
		PasswordHash: "password",
		Role:         identityDomain.Customer,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := userRepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	
	ord := &domain.Order{
		ID:         uuid.New(),
		UserID:     userID,
		Status:     domain.StatusPending,
		TotalPrice: amount,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = repo.CreateOrder(context.Background(), ord)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fetched, err := repo.GetOrderByID(context.Background(), userID, ord.ID)
	if err != nil {
		t.Fatalf("failed to fetch created order: %v", err)
	}

	if fetched.ID != ord.ID {
		t.Errorf("expected ID %v, got %v", ord.ID, fetched.ID)
	}
}

func TestGetOrderByID_NotFound(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	_, err := repo.GetOrderByID(context.Background(), uuid.New(), uuid.New())
	if err != order.ErrOrderNotFound {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}
