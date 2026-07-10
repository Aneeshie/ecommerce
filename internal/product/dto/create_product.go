package dto

import "github.com/google/uuid"

type CreateProductRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Price int64 `json:"price"`
}

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

