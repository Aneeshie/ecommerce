package service

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/google/uuid"
)

type MockUserStore struct {
	FindByEmailFn            func(ctx context.Context, email string) (domain.User, error)
	CreateFn                 func(ctx context.Context, user domain.User) error
	FindByIDFn               func(ctx context.Context, id uuid.UUID) (domain.User, error)
	CreateRefreshTokenFn     func(ctx context.Context, token domain.RefreshToken) error
	FindRefreshTokenByHashFn func(ctx context.Context, hash string) (domain.RefreshToken, error)
}

func (m *MockUserStore) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return domain.User{}, nil
}

func (m *MockUserStore) Create(ctx context.Context, user domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}

func (m *MockUserStore) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return domain.User{}, nil
}

func (m *MockUserStore) CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	if m.CreateRefreshTokenFn != nil {
		return m.CreateRefreshTokenFn(ctx, token)
	}
	return nil
}

func (m *MockUserStore) FindRefreshTokenByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	if m.FindRefreshTokenByHashFn != nil {
		return m.FindRefreshTokenByHashFn(ctx, hash)
	}
	return domain.RefreshToken{}, nil
}
