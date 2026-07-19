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

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) error {
	if req.Email == "" {
		return identity.ErrEmailRequired
	}

	_, err := s.users.FindByEmail(ctx, req.Email)
	if err == nil {
		return identity.ErrEmailAlreadyExists
	}

	if !errors.Is(err, identity.ErrUserNotFound) {
		return err
	}

	err = validator.ValidatePassword(req.Password)
	if err != nil {
		return err
	}

	hash, err := password.Hash(req.Password)

	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	now := time.Now()
	//check if user exists in first place
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, identity.ErrUserNotFound) {
			return dto.LoginResponse{}, identity.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, err
	}

	// check if the password entered is correct
	isValid := password.CompareHash(req.Password, user.PasswordHash)

	if !isValid {
		return dto.LoginResponse{}, identity.ErrInvalidCredentials
	}

	// generate access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user, AccessTokenTTL)

	if err != nil {

		return dto.LoginResponse{}, err
	}

	// generate refresh token
	rawRefreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return dto.LoginResponse{}, err
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
		return dto.LoginResponse{}, err
	}

	if refreshToken.RevokedAt != nil {
		return dto.LoginResponse{}, identity.ErrInvalidRefreshToken
	}
	//return loginResponse

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    int(AccessTokenTTL.Seconds()),
	}, nil

}

func (s *Service) Refresh(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResponse, error) {
	hashedRefreshToken := s.tokenManager.HashRefreshToken(req.RefreshToken)

	refreshToken, err := s.users.FindRefreshTokenByHash(ctx, hashedRefreshToken)

	if refreshToken.RevokedAt != nil {
		return dto.RefreshResponse{}, identity.ErrInvalidRefreshToken
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return dto.RefreshResponse{}, identity.ErrInvalidRefreshToken
	}

	user, err := s.users.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return dto.RefreshResponse{}, err
	}

	// generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user, AccessTokenTTL)
	if err != nil {
		return dto.RefreshResponse{}, err
	}
	//return dto.Refres.usersnse
	return dto.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(AccessTokenTTL.Seconds()),
	}, nil

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
