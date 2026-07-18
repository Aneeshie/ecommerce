package httpx

import (
	"errors"
	"log"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/identity"
)

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, identity.ErrEmailAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, identity.ErrInvalidCredentials):
		http.Error(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, identity.ErrInvalidRefreshToken):
		http.Error(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, identity.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
