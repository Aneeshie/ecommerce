package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/identity/token"
)

type AuthMiddleware struct {
	tokenManager *token.Manager
}

func NewAuthMiddleware(manager *token.Manager) *AuthMiddleware{
	return &AuthMiddleware{
		tokenManager: manager,
	}
}

func (a *AuthMiddleware) Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){

		tokenString, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		customClaims, err := a.tokenManager.VerifyAccessToken(tokenString)
		if err != nil  {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, customClaims)


		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func extractBearerToken(r *http.Request) (string,error) {
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


