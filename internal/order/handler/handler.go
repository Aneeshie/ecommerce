package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	"github.com/Aneeshie/ecommerce/internal/order/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func RegisterRoutes(r chi.Router, h *Handler, auth *middleware.AuthMiddleware) {
	r.With(auth.Auth, auth.RequireRole(domain.Admin)).Post("/api/v1/orders", h.CreateOrder)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "not able to get the jwt claims", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "not able to parse the userID", http.StatusUnauthorized)
		return
	}

	err = h.service.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, "Could not place the order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
