package identity

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrUserNotFound        = errors.New("User not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("Invalid Refresh Token")
	ErrEmailRequired       = errors.New("Email cannot be empty")
)
