package token

import (
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role domain.Role `json:"role"`
	jwt.RegisteredClaims
}
