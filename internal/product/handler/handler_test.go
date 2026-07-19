package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aneeshie/ecommerce/internal/product/dto"
	"github.com/google/uuid"
)

func TestCreateProductHandler(t *testing.T) {
	tests := []struct {
		Name           string
		Payload        interface{}
		SetupMock      func(m *MockProductService)
		ExpectedStatus int
	}{
		{
			Name: "Successful Create",
			Payload: dto.CreateProductRequest{
				Name:        "Test Product",
				Description: "Desc",
				Price:       1000,
			},
			SetupMock: func(m *MockProductService) {
				m.CreateProductFn = func(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
					return &dto.CreateProductResponse{ID: uuid.New()}, nil
				}
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "Invalid Payload",
			Payload:        "invalid-json",
			SetupMock:      func(m *MockProductService) {},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Service Error",
			Payload: dto.CreateProductRequest{
				Name:        "Test Product",
				Description: "Desc",
				Price:       1000,
			},
			SetupMock: func(m *MockProductService) {
				m.CreateProductFn = func(ctx context.Context, req *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
					return nil, errors.New("service error")
				}
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockProductService{}
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateProduct(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}

func TestListProductsHandler(t *testing.T) {
	tests := []struct {
		Name           string
		SetupMock      func(m *MockProductService)
		ExpectedStatus int
	}{
		{
			Name: "Successful List",
			SetupMock: func(m *MockProductService) {
				m.ListProductsFn = func(ctx context.Context, limit int64) ([]*dto.ProductResponse, error) {
					return []*dto.ProductResponse{}, nil
				}
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "Service Error",
			SetupMock: func(m *MockProductService) {
				m.ListProductsFn = func(ctx context.Context, limit int64) ([]*dto.ProductResponse, error) {
					return nil, errors.New("service error")
				}
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockService := &MockProductService{}
			if tt.SetupMock != nil {
				tt.SetupMock(mockService)
			}
			handler := NewHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
			w := httptest.NewRecorder()

			handler.ListProducts(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Fatalf("expected status %d got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
