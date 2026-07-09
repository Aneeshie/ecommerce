package middleware

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/identity/token"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func ClaimsFromContext(ctx context.Context) (*token.CustomClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*token.CustomClaims)
	return claims, ok
}
