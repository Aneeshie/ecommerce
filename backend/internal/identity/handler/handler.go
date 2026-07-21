package handler

import (
	"context"
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

type IdentityService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*service.AuthTokens, error)
	Login(ctx context.Context, req dto.LoginRequest) (*service.AuthTokens, error)
	Refresh(ctx context.Context, refreshTokenString string) (string, error)
	GetCurrentUser(ctx context.Context, userId uuid.UUID) (*dto.MeResponse, error)
}

type Handler struct {
	service      IdentityService
	cookieSecure bool
}

func RegisterRoutes(r chi.Router, h *Handler, auth *middleware.AuthMiddleware) {
	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)
	r.Post("/api/v1/auth/refresh", h.Refresh)

	r.With(auth.Auth).Get("/api/v1/auth/me", h.Me)
}

func NewHandler(s IdentityService, cookieSecure bool) *Handler {
	return &Handler{
		service:      s,
		cookieSecure: cookieSecure,
	}
}

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

	authToken, err := h.service.Register(r.Context(), req)

	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	h.setAuthTokens(w, authToken)

	resp := dto.RegisterResponse{
		Message: "User registered successfully",
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	authToken, err := h.service.Login(r.Context(), req)

	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	h.setAuthTokens(w, authToken)

	httpx.WriteJSON(w, http.StatusOK, "login successful")
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token cookie missing", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	httpx.WriteJSON(w, http.StatusOK, "refresh successful")
}

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

func (h *Handler) setAuthTokens(w http.ResponseWriter, authToken *service.AuthTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    authToken.AccessToken,
		HttpOnly: true,
		Secure:   h.cookieSecure, // false on localhost if not using https
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    authToken.RefreshToken,
		HttpOnly: true,
		Secure:   h.cookieSecure, // false on localhost if not using https
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
