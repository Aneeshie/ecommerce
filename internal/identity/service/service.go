package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/password"
	"github.com/Aneeshie/ecommerce/internal/identity/repository"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
	"github.com/Aneeshie/ecommerce/internal/identity/validator"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("Invalid Refresh Token")
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type Service struct {
	repo         *repository.Repository
	tokenManager *token.Manager
}

func NewService(repo *repository.Repository, tokenManager *token.Manager) *Service {
	return &Service{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("Email cannot be empty")
	}

	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return ErrEmailAlreadyExists
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

	err = s.repo.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	now := time.Now()
	//check if user exists in first place
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Println("find email failure", err)
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	// check if the password entered is correct
	isValid := password.CompareHash(req.Password, user.PasswordHash)

	if !isValid {
		log.Println("password comparison failure")
		return dto.LoginResponse{}, ErrInvalidCredentials
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
	err = s.repo.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return dto.LoginResponse{}, err
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

	refreshToken, err := s.repo.FindRefreshTokenByHash(ctx, hashedRefreshToken)

	if refreshToken.RevokedAt != nil {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	user, err := s.repo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return dto.RefreshResponse{}, err
	}

	// generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user, AccessTokenTTL)
	if err != nil {
		return dto.RefreshResponse{}, err
	}
	//return dto.RefreshReponse
	return dto.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(AccessTokenTTL.Seconds()),
	}, nil

}

func (s *Service) GetCurrentUser(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error) {
	user, err := s.repo.FindByID(ctx, userId)
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
