package store

import (
	identityRepository "github.com/Aneeshie/ecommerce/internal/identity/repository"
	inventoryRepository "github.com/Aneeshie/ecommerce/internal/inventory/repository"
	productRepository "github.com/Aneeshie/ecommerce/internal/product/repository"
	"github.com/jackc/pgx/v5"
)

type TxStore struct {
	tx pgx.Tx
}

func (s *TxStore) Products() *productRepository.Repository {
	return productRepository.NewRepository(s.tx)
}

func (s *TxStore) Inventory() *inventoryRepository.Repository {
	return inventoryRepository.NewRepository(s.tx)
}

func (s *TxStore) Users() *identityRepository.Repository {
	return identityRepository.NewRepository(s.tx)
}
