package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/api/v1/products", h.CreateProduct)
	r.Get("/api/v1/products", h.ListProducts)
	r.Get("/api/v1/products/{id}", h.GetProduct)
	r.Put("/api/v1/products/{id}", h.UpdateProduct)
	r.Delete("/api/v1/products/{id}", h.DeleteProduct)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListProducts(r.Context(), LIMIT)
	if err != nil {
		http.Error(w, "Could not fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to enode response: %v", err)
	}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	resp, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to enode response: %v", err)
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	productID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetProductById(r.Context(), productID)
	if err != nil {
		http.Error(w, "Failed to get the product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to enode response: %v", err)
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	//get the product id from param
	id := chi.URLParam(r, "id")

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
		http.Error(w, "Failed to update the product", http.StatusInternalServerError)
		return
	}

	resp, err := h.service.GetProductById(r.Context(), productID)
	if err != nil {
		http.Error(w, "Failed to get the product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to enode response: %v", err)
	}
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	productID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
