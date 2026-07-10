package domain

import (
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       money.Money
	Status      Status

	CreatedAt time.Time
	UpdatedAt time.Time
}
