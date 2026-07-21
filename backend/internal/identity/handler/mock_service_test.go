package handler

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/google/uuid"
)

type MockIdentityService struct {
	RegisterFn       func(ctx context.Context, req dto.RegisterRequest) (*service.AuthTokens, error)
	LoginFn          func(ctx context.Context, req dto.LoginRequest) (*service.AuthTokens, error)
	RefreshFn        func(ctx context.Context, refreshTokenString string) (string, error)
	GetCurrentUserFn func(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error)
}

func (m *MockIdentityService) Register(ctx context.Context, req dto.RegisterRequest) (*service.AuthTokens, error) {
	if m.RegisterFn != nil {
		return m.RegisterFn(ctx, req)
	}
	return &service.AuthTokens{}, nil
}

func (m *MockIdentityService) Login(ctx context.Context, req dto.LoginRequest) (*service.AuthTokens, error) {
	if m.LoginFn != nil {
		return m.LoginFn(ctx, req)
	}
	return &service.AuthTokens{}, nil
}

func (m *MockIdentityService) Refresh(ctx context.Context, refreshTokenString string) (string, error) {
	if m.RefreshFn != nil {
		return m.RefreshFn(ctx, refreshTokenString)
	}
	return "", nil
}

func (m *MockIdentityService) GetCurrentUser(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error) {
	if m.GetCurrentUserFn != nil {
		return m.GetCurrentUserFn(ctx, userId)
	}
	return nil, nil
}
