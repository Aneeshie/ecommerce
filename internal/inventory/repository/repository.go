package repository

import (
	"context"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/inventory/domain"
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
	_,err := r.db.Exec(ctx, query, inventory.ProductID, inventory.Quantity, inventory.CreatedAt, inventory.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
