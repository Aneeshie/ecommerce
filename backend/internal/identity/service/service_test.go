package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity"
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/password"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	baseReq := dto.RegisterRequest{
		Name:     "Aneesh",
		Email:    "aneesh@example.com",
		Password: "Password123!",
	}

	tests := []struct {
		Name          string
		reqModifier   func(*dto.RegisterRequest)
		setupMock     func(*MockUserStore)
		expectedError error
	}{
		{
			Name: "Successful Registration",
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return domain.User{}, identity.ErrUserNotFound
				}
				m.CreateFn = func(ctx context.Context, user domain.User) error {
					return nil
				}
			},
			expectedError: nil,
		},
		{
			Name: "Email Already Exists",
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return domain.User{}, nil
				}
				m.CreateFn = func(ctx context.Context, user domain.User) error {
					t.Fatal("Create should not have been called.")
					return nil
				}
			},
			expectedError: identity.ErrEmailAlreadyExists,
		},
		{
			Name: "Empty Email",
			reqModifier: func(r *dto.RegisterRequest) {
				r.Email = ""
			},
			setupMock: func(m *MockUserStore) {
				// Should not call any methods
			},
			expectedError: identity.ErrEmailRequired,
		},
		{
			Name: "Database Error on FindByEmail",
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return domain.User{}, errors.New("database error")
				}
				m.CreateFn = func(ctx context.Context, user domain.User) error {
					t.Fatal("Create should not have been called.")
					return nil
				}
			},
			expectedError: errors.New("database error"),
		},
		{
			Name: "Database Error on Create",
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return domain.User{}, identity.ErrUserNotFound
				}
				m.CreateFn = func(ctx context.Context, user domain.User) error {
					return errors.New("insert error")
				}
			},
			expectedError: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mock := &MockUserStore{}

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			service := &Service{
				users:        mock,
				tokenManager: token.NewManager("test"),
			}

			req := baseReq
			if tt.reqModifier != nil {
				tt.reqModifier(&req)
			}

			//TODO
			tokens, err := service.Register(context.Background(), req)

			// We use string comparison for dynamic errors like "database error" if they don't match exactly by pointer
			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				if tokens == nil {
					t.Fatalf("expected tokens to not be nil on success")
				}
				if tokens.AccessToken == "" || tokens.RefreshToken == "" {
					t.Fatalf("expected tokens to be populated, got access: %q, refresh: %q", tokens.AccessToken, tokens.RefreshToken)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	hashedPassword, _ := password.Hash("Password123!")
	validUser := domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	tests := []struct {
		Name          string
		req           dto.LoginRequest
		setupMock     func(*MockUserStore)
		expectedError error
	}{
		{
			Name: "Successful Login",
			req:  dto.LoginRequest{Email: "test@example.com", Password: "Password123!"},
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return validUser, nil
				}
				m.CreateRefreshTokenFn = func(ctx context.Context, token domain.RefreshToken) error {
					return nil
				}
			},
			expectedError: nil,
		},
		{
			Name: "User Not Found",
			req:  dto.LoginRequest{Email: "notfound@example.com", Password: "Password123!"},
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return domain.User{}, identity.ErrUserNotFound
				}
			},
			expectedError: identity.ErrInvalidCredentials,
		},
		{
			Name: "Invalid Password",
			req:  dto.LoginRequest{Email: "test@example.com", Password: "WrongPassword!"},
			setupMock: func(m *MockUserStore) {
				m.FindByEmailFn = func(ctx context.Context, email string) (domain.User, error) {
					return validUser, nil
				}
			},
			expectedError: identity.ErrInvalidCredentials,
		},
	}

	tm := token.NewManager("test-secret")

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mock := &MockUserStore{}
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			service := &Service{
				users:        mock,
				tokenManager: tm,
			}

			tokens, err := service.Login(context.Background(), tt.req)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				if tokens == nil {
					t.Fatalf("expected tokens to not be nil on success")
				}
				if tokens.AccessToken == "" || tokens.RefreshToken == "" {
					t.Fatalf("expected tokens to be populated, got access: %q, refresh: %q", tokens.AccessToken, tokens.RefreshToken)
				}
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	validUser := domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tm := token.NewManager("test-secret")
	now := time.Now()

	tests := []struct {
		Name          string
		req           string
		setupMock     func(*MockUserStore)
		expectedError error
	}{
		{
			Name: "Successful Refresh",
			req:  "valid-token",
			setupMock: func(m *MockUserStore) {
				m.FindRefreshTokenByHashFn = func(ctx context.Context, hash string) (domain.RefreshToken, error) {
					return domain.RefreshToken{
						UserID:    validUser.ID,
						ExpiresAt: now.Add(1 * time.Hour),
					}, nil
				}
				m.FindByIDFn = func(ctx context.Context, id uuid.UUID) (domain.User, error) {
					return validUser, nil
				}
			},
			expectedError: nil,
		},
		{
			Name: "Expired Token",
			req:  "expired-token",
			setupMock: func(m *MockUserStore) {
				m.FindRefreshTokenByHashFn = func(ctx context.Context, hash string) (domain.RefreshToken, error) {
					return domain.RefreshToken{
						UserID:    validUser.ID,
						ExpiresAt: now.Add(-1 * time.Hour), // Expired
					}, nil
				}
			},
			expectedError: identity.ErrInvalidRefreshToken,
		},
		{
			Name: "Revoked Token",
			req:  "revoked-token",
			setupMock: func(m *MockUserStore) {
				m.FindRefreshTokenByHashFn = func(ctx context.Context, hash string) (domain.RefreshToken, error) {
					revokedAt := now.Add(-1 * time.Hour)
					return domain.RefreshToken{
						UserID:    validUser.ID,
						ExpiresAt: now.Add(1 * time.Hour),
						RevokedAt: &revokedAt, // Revoked
					}, nil
				}
			},
			expectedError: identity.ErrInvalidRefreshToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mock := &MockUserStore{}
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			service := &Service{
				users:        mock,
				tokenManager: tm,
			}

			accessToken, err := service.Refresh(context.Background(), tt.req)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				if accessToken == "" {
					t.Fatalf("expected access token to not be empty")
				}
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	validID := uuid.New()
	validUser := domain.User{
		ID:    validID,
		Name:  "Aneesh",
		Email: "test@example.com",
		Role:  domain.Customer,
	}

	tests := []struct {
		Name          string
		userID        uuid.UUID
		setupMock     func(*MockUserStore)
		expectedError error
		expectedEmail string
	}{
		{
			Name:   "Successful Get",
			userID: validID,
			setupMock: func(m *MockUserStore) {
				m.FindByIDFn = func(ctx context.Context, id uuid.UUID) (domain.User, error) {
					return validUser, nil
				}
			},
			expectedError: nil,
			expectedEmail: validUser.Email,
		},
		{
			Name:   "User Not Found",
			userID: uuid.New(),
			setupMock: func(m *MockUserStore) {
				m.FindByIDFn = func(ctx context.Context, id uuid.UUID) (domain.User, error) {
					return domain.User{}, identity.ErrUserNotFound
				}
			},
			expectedError: identity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mock := &MockUserStore{}
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			service := &Service{
				users: mock,
			}

			res, err := service.GetCurrentUser(context.Background(), tt.userID)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				if res.Email != tt.expectedEmail {
					t.Fatalf("expected email %v got %v", tt.expectedEmail, res.Email)
				}
			}
		})
	}
}
