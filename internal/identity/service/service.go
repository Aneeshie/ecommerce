package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/password"
	"github.com/Aneeshie/ecommerce/internal/identity/repository"
	"github.com/Aneeshie/ecommerce/internal/identity/validator"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service{
	return &Service{
		repo: repo,
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
