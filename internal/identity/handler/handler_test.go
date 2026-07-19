package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aneeshie/ecommerce/internal/identity/dto"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		Name           string
		Payload        interface{}
		SetupMock      func(m *MockIdentityService)
		ExpectedStatus int
	}{
		{
			Name: "Successful Registration",
			Payload: dto.RegisterRequest{
				Name:     "Aneesh",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			SetupMock: func(m *MockIdentityService) {
				m.RegisterFn = func(ctx context.Context, req dto.RegisterRequest) error {
					return nil
				}
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name: "Invalid Payload",
			Payload: "invalid-json-string",
			SetupMock: func(m *MockIdentityService) {
				// Should not be called
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Service Error",
			Payload: dto.RegisterRequest{
				Name:     "Aneesh",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			SetupMock: func(m *MockIdentityService) {
				m.RegisterFn = func(ctx context.Context, req dto.RegisterRequest) error {
					return errors.New("some error")
				}
			},
			ExpectedStatus: http.StatusInternalServerError, // httpx.WriteError returns 500 for generic errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockIdentityService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			var b []byte
			if s, ok := tt.Payload.(string); ok {
				b = []byte(s)
			} else {
				b, _ = json.Marshal(tt.Payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		Name           string
		Payload        interface{}
		SetupMock      func(m *MockIdentityService)
		ExpectedStatus int
	}{
		{
			Name: "Successful Login",
			Payload: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			SetupMock: func(m *MockIdentityService) {
				m.LoginFn = func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
					return dto.LoginResponse{AccessToken: "token"}, nil
				}
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "Invalid Payload",
			Payload: "invalid-json-string",
			SetupMock: func(m *MockIdentityService) {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Service Error",
			Payload: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			SetupMock: func(m *MockIdentityService) {
				m.LoginFn = func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
					return dto.LoginResponse{}, errors.New("invalid credentials")
				}
			},
			ExpectedStatus: http.StatusInternalServerError, // httpx.WriteError defaults to 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockIdentityService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			var b []byte
			if s, ok := tt.Payload.(string); ok {
				b = []byte(s)
			} else {
				b, _ = json.Marshal(tt.Payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		Name           string
		Payload        interface{}
		SetupMock      func(m *MockIdentityService)
		ExpectedStatus int
	}{
		{
			Name: "Successful Refresh",
			Payload: dto.RefreshRequest{
				RefreshToken: "some-token",
			},
			SetupMock: func(m *MockIdentityService) {
				m.RefreshFn = func(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResponse, error) {
					return dto.RefreshResponse{AccessToken: "new-token"}, nil
				}
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "Invalid Payload",
			Payload: "invalid-json-string",
			SetupMock: func(m *MockIdentityService) {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Service Error",
			Payload: dto.RefreshRequest{
				RefreshToken: "some-token",
			},
			SetupMock: func(m *MockIdentityService) {
				m.RefreshFn = func(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResponse, error) {
					return dto.RefreshResponse{}, errors.New("invalid refresh token")
				}
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockIdentityService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			var b []byte
			if s, ok := tt.Payload.(string); ok {
				b = []byte(s)
			} else {
				b, _ = json.Marshal(tt.Payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Refresh(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}

func TestMeHandler(t *testing.T) {
	tests := []struct {
		Name           string
		SetupMock      func(m *MockIdentityService)
		ExpectedStatus int
	}{
		{
			Name: "Unauthorized (No Claims in Context)",
			SetupMock: func(m *MockIdentityService) {
			},
			ExpectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockIdentityService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			w := httptest.NewRecorder()

			handler.Me(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
