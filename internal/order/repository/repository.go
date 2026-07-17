package repository

import (
	"context"
	"errors"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/google/uuid"
)

type Repository struct {
	db database.QueryExecutor
}

func NewRepository(db database.QueryExecutor) *Repository {
	return &Repository{
		db: db,
	}
}

var ErrInsufficientInventory error = errors.New("Insufficient inventory")

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

func (r *Repository) CreateOrderItems(ctx context.Context, orderItems []*domain.OrderItem) error {

	query := `INSERT INTO order_items (
    id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7);`

	for _, item := range orderItems {
		_, err := r.db.Exec(ctx, query, item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price.Amount(), item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DecreaseInventory(ctx context.Context, productID uuid.UUID, quantity int) error {
	query := `UPDATE inventories SET quantity = quantity - $1 WHERE product_id = $2 AND quantity >= $1`

	result, err := r.db.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return ErrInsufficientInventory
	}

	return nil
}
