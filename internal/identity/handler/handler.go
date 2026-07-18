package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Aneeshie/ecommerce/internal/httpx"
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
	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)
	r.Post("/api/v1/auth/refresh", h.Refresh)

	r.With(auth.Auth).Get("/auth/me", h.Me)
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

// Register godoc
//
//	@Summary Register a new user
//	@Description Creates a new customer account
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Param request body dto.RegisterRequest true "Register Request"
//	@Success 201 {object} dto.RegisterResponse
//	@Failure 400 {string} string
//	@Failure 409 {string} string
//	@Failure 500 {string} string
//	@Router /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	err = h.service.Register(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	resp := dto.RegisterResponse{
		Message: "User registered successfully",
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)

}

// Login godoc
//
//	@Summary Login user
//	@Description Authenticates a user and returns tokens
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Param request body dto.LoginRequest true "Login Request"
//	@Success 200 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 500 {string} string
//	@Router /api/v1/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	resp, err := h.service.Login(r.Context(), req)

	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Refresh godoc
//
//	@Summary Refresh access token
//	@Description Refreshes the authentication token
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Param request body dto.RefreshRequest true "Refresh Request"
//	@Success 200 {object} map[string]interface{}
//	@Failure 400 {string} string
//	@Failure 401 {string} string
//	@Failure 500 {string} string
//	@Router /api/v1/auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Refresh(r.Context(), req)

	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Me godoc
//
//	@Summary Get current user
//	@Description Returns the profile of the currently authenticated user
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Success 200 {object} map[string]interface{}
//	@Failure 401 {string} string
//	@Failure 500 {string} string
//	@Security ApiKeyAuth
//	@Router /auth/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)
}
