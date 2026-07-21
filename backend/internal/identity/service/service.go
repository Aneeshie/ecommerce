package service

import (
	"context"
	"errors"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity"
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/password"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
	"github.com/Aneeshie/ecommerce/internal/identity/validator"
	"github.com/Aneeshie/ecommerce/internal/store"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type UserStore interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, hash string) (domain.RefreshToken, error)
	Create(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type Service struct {
	users        UserStore
	tokenManager *token.Manager
}

func NewService(store *store.Store, tokenManager *token.Manager) *Service {
	return &Service{
		users:        store.Users(),
		tokenManager: tokenManager,
	}
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) (*AuthTokens, error) {
	if req.Email == "" {
		return nil, identity.ErrEmailRequired
	}

	_, err := s.users.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, identity.ErrEmailAlreadyExists
	}

	if !errors.Is(err, identity.ErrUserNotFound) {
		return nil, err
	}

	err = validator.ValidatePassword(req.Password)
	if err != nil {
		return nil, err
	}

	hash, err := password.Hash(req.Password)

	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:            uuid.New(),
		Name:          req.Name,
		Email:         req.Email,
		PasswordHash:  string(hash),
		Role:          domain.Customer,
		EmailVerified: false, //TODO: EMAIL VERIFICATION
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.users.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, &user)

}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (*AuthTokens, error) {
	//check if user exists in first place
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, identity.ErrUserNotFound) {
			return nil, identity.ErrInvalidCredentials

		}
		return nil, err
	}

	// check if the password entered is correct
	isValid := password.CompareHash(req.Password, user.PasswordHash)

	if !isValid {
		return nil, identity.ErrInvalidCredentials
	}

	return s.generateTokens(ctx, &user)

}

func (s *Service) Refresh(ctx context.Context, refreshTokenString string) (string, error) {
	hashedRefreshToken := s.tokenManager.HashRefreshToken(refreshTokenString)

	refreshToken, err := s.users.FindRefreshTokenByHash(ctx, hashedRefreshToken)

	if refreshToken.RevokedAt != nil {
		return "", identity.ErrInvalidRefreshToken
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return "", identity.ErrInvalidRefreshToken
	}

	user, err := s.users.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return "", err
	}

	// generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user, AccessTokenTTL)
	if err != nil {
		return "", err
	}
	
	return accessToken, nil

}

func (s *Service) GetCurrentUser(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error) {
	user, err := s.users.FindByID(ctx, userId)
	if err != nil {
		return &dto.MeResponse{}, err
	}

	return &dto.MeResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *Service) generateTokens(ctx context.Context, user *domain.User) (*AuthTokens, error) {
	now := time.Now()
	// generate access token
	accessToken, err := s.tokenManager.GenerateAccessToken(*user, AccessTokenTTL)

	if err != nil {

		return nil, err
	}

	// generate refresh token
	rawRefreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// hash refresh token
	hashedRefreshToken := s.tokenManager.HashRefreshToken(rawRefreshToken)

	// create refreshToken entity
	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: now.Add(RefreshTokenTTL),
		CreatedAt: now,
		RevokedAt: nil,
	}

	//call the repo thingy (put it in the database)
	err = s.users.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// if refreshToken.RevokedAt != nil {
	// 	return nil, identity.ErrInvalidRefreshToken
	// }

	return &AuthTokens{AccessToken: accessToken, RefreshToken: rawRefreshToken}, nil
}
