package handler

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/google/uuid"
)

type MockIdentityService struct {
	RegisterFn       func(ctx context.Context, req dto.RegisterRequest) error
	LoginFn          func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	RefreshFn        func(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResponse, error)
	GetCurrentUserFn func(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error)
}

func (m *MockIdentityService) Register(ctx context.Context, req dto.RegisterRequest) error {
	if m.RegisterFn != nil {
		return m.RegisterFn(ctx, req)
	}
	return nil
}

func (m *MockIdentityService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	if m.LoginFn != nil {
		return m.LoginFn(ctx, req)
	}
	return dto.LoginResponse{}, nil
}

func (m *MockIdentityService) Refresh(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResponse, error) {
	if m.RefreshFn != nil {
		return m.RefreshFn(ctx, req)
	}
	return dto.RefreshResponse{}, nil
}

func (m *MockIdentityService) GetCurrentUser(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error) {
	if m.GetCurrentUserFn != nil {
		return m.GetCurrentUserFn(ctx, userId)
	}
	return nil, nil
}
