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
	"github.com/Aneeshie/ecommerce/internal/identity/service"
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
				m.RegisterFn = func(ctx context.Context, req dto.RegisterRequest) (*service.AuthTokens, error) {
					return &service.AuthTokens{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil
				}
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:    "Invalid Payload",
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
				m.RegisterFn = func(ctx context.Context, req dto.RegisterRequest) (*service.AuthTokens, error) {
					return nil, errors.New("some error")
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
			handler := NewHandler(mockService, false)

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

			if w.Code == http.StatusCreated {
				cookies := w.Result().Cookies()
				if len(cookies) != 2 {
					t.Fatalf("expected 2 cookies, got %d", len(cookies))
				}
				for _, c := range cookies {
					if !c.HttpOnly {
						t.Errorf("expected cookie %s to be HttpOnly", c.Name)
					}
				}
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
				m.LoginFn = func(ctx context.Context, req dto.LoginRequest) (*service.AuthTokens, error) {
					return &service.AuthTokens{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil
				}
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:    "Invalid Payload",
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
				m.LoginFn = func(ctx context.Context, req dto.LoginRequest) (*service.AuthTokens, error) {
					return &service.AuthTokens{}, errors.New("invalid credentials")
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
			handler := NewHandler(mockService, false)

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

			if w.Code == http.StatusOK {
				cookies := w.Result().Cookies()
				if len(cookies) != 2 {
					t.Fatalf("expected 2 cookies, got %d", len(cookies))
				}
				for _, c := range cookies {
					if !c.HttpOnly {
						t.Errorf("expected cookie %s to be HttpOnly", c.Name)
					}
				}
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		Name           string
		CookieValue    string
		SetupMock      func(m *MockIdentityService)
		ExpectedStatus int
	}{
		{
			Name:        "Successful Refresh",
			CookieValue: "some-token",
			SetupMock: func(m *MockIdentityService) {
				m.RefreshFn = func(ctx context.Context, refreshTokenString string) (string, error) {
					return "new-token", nil
				}
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:        "Missing Cookie",
			CookieValue: "",
			SetupMock: func(m *MockIdentityService) {
			},
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name:        "Service Error",
			CookieValue: "invalid-token",
			SetupMock: func(m *MockIdentityService) {
				m.RefreshFn = func(ctx context.Context, refreshTokenString string) (string, error) {
					return "", errors.New("invalid refresh token")
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
			handler := NewHandler(mockService, false)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
			if tt.CookieValue != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.CookieValue})
			}
			w := httptest.NewRecorder()

			handler.Refresh(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				cookies := w.Result().Cookies()
				if len(cookies) != 1 {
					t.Fatalf("expected 1 cookie (access_token), got %d", len(cookies))
				}
				if cookies[0].Name != "access_token" || !cookies[0].HttpOnly {
					t.Errorf("expected HttpOnly access_token cookie")
				}
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
			handler := NewHandler(mockService, false)

			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			w := httptest.NewRecorder()

			handler.Me(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
