package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ProductID uuid.UUID
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
