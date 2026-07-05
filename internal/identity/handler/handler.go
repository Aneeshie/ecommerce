package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
}

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/auth/register", h.Register)
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	err = h.service.Register(r.Context(), req)

	if err != nil {
    log.Printf("register error: %+v\n", err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
	}

	if errors.Is(err, service.ErrEmailAlreadyExists) {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}


	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	resp := dto.RegisterResponse{
		Message: "User registered successfully",
	}
	json.NewEncoder(w).Encode(resp)
}
