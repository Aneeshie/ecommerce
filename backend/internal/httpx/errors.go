package httpx

import (
	"errors"
	"log"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/identity"
	"github.com/Aneeshie/ecommerce/internal/product"
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

	case errors.Is(err, product.ErrProductNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, product.ErrEmptyProductName):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, product.ErrEmptyProductDescription):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, money.ErrNegativeAmount):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
