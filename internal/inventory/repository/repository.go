package repository

import (
	"github.com/Aneeshie/ecommerce/internal/common/database"
)

type Repository struct {
	db database.QueryExecutor
}

func NewRepository(db database.QueryExecutor) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateInventory() {

}
