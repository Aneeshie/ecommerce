package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/httpx"
	identityDomain "github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error)
}

type Handler struct {
	service OrderService
}

func NewHandler(service OrderService) *Handler {
	return &Handler{
		service: service,
	}
}

func RegisterRoutes(r chi.Router, h *Handler, auth *middleware.AuthMiddleware) {
	r.With(auth.Auth, auth.RequireRole(identityDomain.Admin)).Post("/api/v1/orders", h.CreateOrder)
	r.Get("/api/v1/orders", h.GetOrders)
	r.Get("/api/v1/orders/{orderID}", h.GetOrderByID)
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
		httpx.WriteError(w, err)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	err = h.service.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, "order created successfully")

}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	orders, err := h.service.GetOrdersByUserID(
		r.Context(),
		userID,
	)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, orders)
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(
		r.Context(),
		userID,
		orderID,
	)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}
