package dto

import (
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/google/uuid"
)

type MeResponse struct {
    ID    uuid.UUID   `json:"id"`
    Name  string      `json:"name"`
    Email string      `json:"email"`
    Role  domain.Role `json:"role"`
}
