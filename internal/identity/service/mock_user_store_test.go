package service

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/google/uuid"
)

type MockUserStore struct {
	FindByEmailFn func(ctx context.Context, email string) (domain.User, error)
	CreateFn      func(ctx context.Context, user domain.User) error
}

func (m *MockUserStore) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.FindByEmailFn(ctx, email)
}

func (m *MockUserStore) Create(ctx context.Context, user domain.User) error {
	return m.CreateFn(ctx, user)
}

func (m *MockUserStore) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	panic("not implemented")
}

func (m *MockUserStore) CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	panic("not implemented")
}

func (m *MockUserStore) FindRefreshTokenByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	panic("not implemented")
}
