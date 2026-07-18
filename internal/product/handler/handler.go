package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/httpx"
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/Aneeshie/ecommerce/internal/product/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const LIMIT = 20

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func RegisterRoutes(r chi.Router, h *Handler, auth *middleware.AuthMiddleware) {
	r.With(auth.Auth, auth.RequireRole(domain.Admin)).Post("/api/v1/products", h.CreateProduct)
	r.Get("/api/v1/products", h.ListProducts)
	r.Get("/api/v1/products/{id}", h.GetProduct)
	r.With(auth.Auth, auth.RequireRole(domain.Admin)).Put("/api/v1/products/{id}", h.UpdateProduct)
	r.With(auth.Auth, auth.RequireRole(domain.Admin)).Delete("/api/v1/products/{id}", h.DeleteProduct)
}

// ListProducts godoc
//
//	@Summary List products
//	@Description Returns a list of products
//	@Tags Products
//	@Accept json
//	@Produce json
//	@Success 200 {array} map[string]interface{}
//	@Failure 500 {string} string
//	@Router /api/v1/products [get]
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListProducts(r.Context(), LIMIT)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// CreateProduct godoc
//
//	@Summary Create product
//	@Description Creates a new product (Admin only)
//	@Tags Products
//	@Accept json
//	@Produce json
//	@Param request body dto.CreateProductRequest true "Create Product Request"
//	@Success 201 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 403 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	resp, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

// GetProduct godoc
//
//	@Summary Get product by ID
//	@Description Returns a specific product by its ID
//	@Tags Products
//	@Accept json
//	@Produce json
//	@Param id path string true "Product ID"
//	@Success 200 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 404 {string} string
//	@Failure 500 {string} string
//	@Router /api/v1/products/{id} [get]
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	productID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetProductById(r.Context(), productID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// UpdateProduct godoc
//
//	@Summary Update product
//	@Description Updates a specific product (Admin only)
//	@Tags Products
//	@Accept json
//	@Produce json
//	@Param id path string true "Product ID"
//	@Param request body dto.UpdateProductRequest true "Update Product Request"
//	@Success 200 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 403 {string} string
//	@Failure 404 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/products/{id} [put]
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	//get the product id from param
	id := chi.URLParam(r, "id")

	defer r.Body.Close()

	productID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//modify the product
	err = h.service.UpdateProduct(r.Context(), productID, &req)

	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	resp, err := h.service.GetProductById(r.Context(), productID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// DeleteProduct godoc
//
//	@Summary Delete product
//	@Description Deletes a specific product (Admin only)
//	@Tags Products
//	@Accept json
//	@Produce json
//	@Param id path string true "Product ID"
//	@Success 204 {string} string "Product successfully deleted"
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 403 {string} string
//	@Failure 404 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /api/v1/products/{id} [delete]
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	productID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, "Product successfully deleted")
}
