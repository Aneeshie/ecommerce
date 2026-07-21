package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/testdb"
	"github.com/Aneeshie/ecommerce/internal/identity"
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	now := time.Now()
	user := domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "password",
		Role:         domain.Customer,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Try fetching the created user
	fetchedUser, err := repo.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("failed to fetch created user: %v", err)
	}

	if fetchedUser.ID != user.ID {
		t.Errorf("expected ID %v, got %v", user.ID, fetchedUser.ID)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	tx := testdb.SetupTestDB(t)
	repo := NewRepository(tx)

	_, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
	if err != identity.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
