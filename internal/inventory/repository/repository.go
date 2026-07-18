package repository

import (
	"context"
	"errors"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/inventory"
	"github.com/Aneeshie/ecommerce/internal/inventory/domain"
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

func (r *Repository) CreateInventory(ctx context.Context, inventory *domain.Inventory) error {
	query := `INSERT INTO inventories (product_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, inventory.ProductID, inventory.Quantity, inventory.CreatedAt, inventory.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetInventoryByProductID(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error) {

	var inv domain.Inventory

	query := `
	SELECT
		product_id,
		quantity,
		created_at,
		updated_at
	FROM inventories
	WHERE product_id = $1
`

	err := r.db.QueryRow(ctx, query, productID).Scan(
		&inv.ProductID,
		&inv.Quantity,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, inventory.ErrNoProductFound
		}
	}

	return &inv, nil

}
