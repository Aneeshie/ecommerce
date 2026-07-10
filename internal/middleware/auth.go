package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
)

type AuthMiddleware struct {
	tokenManager *token.Manager
}

func NewAuthMiddleware(manager *token.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: manager,
	}
}

func (a *AuthMiddleware) Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		customClaims, err := a.tokenManager.VerifyAccessToken(tokenString)
		if err != nil {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, customClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("empty jwt token")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid format")
	}

	return parts[1], nil
}

func (a *AuthMiddleware) RequireRole(role domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
