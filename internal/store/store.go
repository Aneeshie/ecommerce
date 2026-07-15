package store

import (
	"context"

	identityRepository "github.com/Aneeshie/ecommerce/internal/identity/repository"
	inventoryRepository "github.com/Aneeshie/ecommerce/internal/inventory/repository"
	productRepository "github.com/Aneeshie/ecommerce/internal/product/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool: pool,
	}
}

func (s *Store) Begin(ctx context.Context) (*TxStore, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return &TxStore{}, err
	}

	return &TxStore{
		tx: tx,
	}, nil
}

func (s *Store) Products() *productRepository.Repository {
	return productRepository.NewRepository(s.pool)
}

func (s *Store) Inventory() *inventoryRepository.Repository {
	return inventoryRepository.NewRepository(s.pool)
}

func (s *Store) Users() *identityRepository.Repository {
	return identityRepository.NewRepository(s.pool)
}
