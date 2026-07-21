package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

		tokenString, err := extractAccessToken(r)
		if err != nil {
			http.Error(w, "Access token cookie missing", http.StatusUnauthorized)
			return
		}

		customClaims, err := a.tokenManager.VerifyAccessToken(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, customClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func extractAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", fmt.Errorf("access token cookie not found: %v", err)
	}

	return cookie.Value, nil
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
				log.Println("Role from context:", role)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
