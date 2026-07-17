package domain

import (
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID

	Quantity int
	Price    money.Money

	CreatedAt time.Time
	UpdatedAt time.Time
}
