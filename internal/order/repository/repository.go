package repository

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
)

type Repository struct {
	db database.QueryExecutor
}

func NewRepository(db database.QueryExecutor) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateOrder(ctx context.Context, order *domain.Order) error {

	query := `
		INSERT INTO orders (
			id,
			user_id,
			status,
			total_price,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, order.ID, order.UserID, order.Status, order.TotalPrice.Amount(), order.CreatedAt, order.UpdatedAt)

	return err
}
