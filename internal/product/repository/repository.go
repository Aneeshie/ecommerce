package repository

import (
	"context"
	"errors"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	"github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price.Amount(), product.Status, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListProducts(ctx context.Context, limit int64) ([]*domain.Product, error) {
	query := `SELECT id, name, description, price, status, created_at, updated_at FROM products LIMIT $1;`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*domain.Product

	for rows.Next() {
		var p domain.Product

		var amount int64

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &amount, &p.Status, &p.CreatedAt, &p.UpdatedAt)

		if err != nil {

			return nil, err
		}

		p.Price, err = money.New(amount)

		if err != nil {
			return nil, err
		}

		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) GetProductByID(ctx context.Context, productId uuid.UUID) (*domain.Product, error) {

	query := `SELECT id, name, description, price, status, created_at, updated_at FROM products WHERE id=$1;`

	var p domain.Product

	var amount int64

	err := r.db.QueryRow(ctx, query, productId).Scan(&p.ID, &p.Name, &p.Description, &amount, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	p.Price, err  = money.New(amount)

	if err != nil {
			return nil, err
		}



	return &p, nil
}
