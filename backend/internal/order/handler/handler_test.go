package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Aneeshie/ecommerce/internal/identity/token"
	"github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	"github.com/google/uuid"
)

func TestCreateOrderHandler(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		Name           string
		Payload        interface{}
		Claims         *token.CustomClaims
		SetupMock      func(m *MockOrderService)
		ExpectedStatus int
	}{
		{
			Name: "Successful Create",
			Payload: dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: uuid.New(), Quantity: 1},
				},
			},
			Claims: &token.CustomClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: userID.String(),
				},
			},
			SetupMock: func(m *MockOrderService) {
				m.CreateOrderFn = func(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error {
					return nil
				}
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "Invalid Payload",
			Payload:        "invalid-json",
			Claims:         nil,
			SetupMock:      func(m *MockOrderService) {},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Missing Claims",
			Payload: dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: uuid.New(), Quantity: 1},
				},
			},
			Claims:         nil, // Missing claims
			SetupMock:      func(m *MockOrderService) {},
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			Name: "Service Error",
			Payload: dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: uuid.New(), Quantity: 1},
				},
			},
			Claims: &token.CustomClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: userID.String(),
				},
			},
			SetupMock: func(m *MockOrderService) {
				m.CreateOrderFn = func(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error {
					return errors.New("service error")
				}
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockOrderService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			var b []byte
			if s, ok := tt.Payload.(string); ok {
				b = []byte(s)
			} else {
				b, _ = json.Marshal(tt.Payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			if tt.Claims != nil {
				req = req.WithContext(middleware.WithClaims(req.Context(), tt.Claims))
			}

			w := httptest.NewRecorder()

			handler.CreateOrder(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
