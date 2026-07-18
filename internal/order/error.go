package order

import "errors"

var (
	ErrOrderNotFound         = errors.New("Order not found.")
	ErrEmptyOrder            = errors.New("The order is empty")
	ErrInvalidQuantity       = errors.New("One of the items has invalid quantity.")
	ErrDuplicateProduct      = errors.New("There are duplicate products.")
	ErrInsufficientInventory = errors.New("Insufficient quantity")
	ErrUnauthorized          = errors.New("unauthorized user")
)
