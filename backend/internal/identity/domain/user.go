package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	PasswordHash  string
	Role          Role
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Role string

const (
	Customer Role = "CUSTOMER"
	Admin    Role = "ADMIN"
)
