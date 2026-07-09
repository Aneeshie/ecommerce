package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/identity/dto"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func RegisterRoutes(r chi.Router, h *Handler, auth *middleware.AuthMiddleware) {
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/refresh", h.Refresh)

	r.With(auth.Auth).Get("/auth/me", h.Me)
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	resp, err := h.service.Login(r.Context(), req)

	if errors.Is(err, service.ErrInvalidCredentials) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("login failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode login response: %v", err)
	}
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Refresh(r.Context(), req)

	if errors.Is(err, service.ErrInvalidRefreshToken) {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("refresh failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode refresh response: %v", err)
	}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Could not parse user Id from claims: %v ", err)
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Failed to get the current user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("failed to encode refresh response: %v", err)
	}

}
