package domain

import (
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     OrderStatus
	TotalPrice money.Money

	CreatedAt time.Time
	UpdatedAt time.Time
}
