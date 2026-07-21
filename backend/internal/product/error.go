package product

import "errors"

var (
	ErrProductNotFound         = errors.New("Product Not Found")
	ErrEmptyProductName        = errors.New("Product name cannot be empty")
	ErrEmptyProductDescription = errors.New("Product description cannot be empty")
)
