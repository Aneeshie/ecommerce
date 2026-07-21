package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/order"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
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
			return fmt.Errorf("create order item: %w", err)
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
		return order.ErrInsufficientInventory
	}

	return nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	query := `
		SELECT
			id,
			user_id,
			status,
			total_price,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		order := &domain.Order{}

		var totalPrice int64

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&totalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		order.TotalPrice, err = money.New(totalPrice)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetOrderByID(
	ctx context.Context,
	userID uuid.UUID,
	orderID uuid.UUID,
) (*domain.Order, error) {

	query := `
		SELECT
			id,
			user_id,
			status,
			total_price,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		  AND user_id = $2
	`

	ord := &domain.Order{}

	var totalPrice int64

	err := r.db.QueryRow(ctx, query, orderID, userID).Scan(
		&ord.ID,
		&ord.UserID,
		&ord.Status,
		&totalPrice,
		&ord.CreatedAt,
		&ord.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order.ErrOrderNotFound
		}
	}

	ord.TotalPrice, err = money.New(totalPrice)
	if err != nil {
		return nil, err
	}

	return ord, nil
}
