package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/Aneeshie/ecommerce/internal/product/service"
	"github.com/go-chi/chi/v5"
)

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
