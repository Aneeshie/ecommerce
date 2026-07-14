package domain

import "time"

type Inventory struct {
	ProductID string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
