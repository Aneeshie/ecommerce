package handler

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	"github.com/google/uuid"
)

type MockOrderService struct {
	CreateOrderFn       func(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error
	GetOrdersByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetOrderByIDFn      func(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error)
}

func (m *MockOrderService) CreateOrder(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error {
	if m.CreateOrderFn != nil {
		return m.CreateOrderFn(ctx, userID, req)
	}
	return nil
}

func (m *MockOrderService) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	if m.GetOrdersByUserIDFn != nil {
		return m.GetOrdersByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockOrderService) GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error) {
	if m.GetOrderByIDFn != nil {
		return m.GetOrderByIDFn(ctx, userID, orderID)
	}
	return nil, nil
}
