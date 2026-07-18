package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/httpx"
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
	r.Get("/api/v1/orders", h.GetOrders)
	r.Get("/api/v1/orders/{orderID}", h.GetOrderByID)
}

// CreateOrder godoc
//
//	@Summary Create order
//	@Description Creates a new order for the authenticated user
//	@Tags Orders
//	@Accept json
//	@Produce json
//	@Param request body dto.CreateOrderRequest true "Create Order Request"
//	@Success 201 {string} string "order created successfully"
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/orders [post]
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

// GetOrders godoc
//
//	@Summary Get user orders
//	@Description Returns a list of orders for the authenticated user
//	@Tags Orders
//	@Accept json
//	@Produce json
//	@Success 200 {array} map[string]interface{}
//	@Failure 401 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/orders [get]
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

// GetOrderByID godoc
//
//	@Summary Get order by ID
//	@Description Returns a specific order by its ID
//	@Tags Orders
//	@Accept json
//	@Produce json
//	@Param orderID path string true "Order ID"
//	@Success 200 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 404 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/orders/{orderID} [get]
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
