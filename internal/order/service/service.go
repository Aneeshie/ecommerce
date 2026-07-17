package service

import (
	"context"
	"errors"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	productDomain "github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/Aneeshie/ecommerce/internal/store"
	"github.com/google/uuid"
)

var (
	ErrEmptyOrder            = errors.New("The order is empty")
	ErrInvalidQuantity       = errors.New("One of the items has invalid quantity.")
	ErrDuplicateProduct      = errors.New("There are duplicate products.")
	ErrInsufficientInventory = errors.New("Insufficient quantity")
	ErrUnauthorized          = errors.New("unauthorized user")
)

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error {
	// validate the request
	if len(req.Items) == 0 {
		return ErrEmptyOrder
	}
	seen := make(map[uuid.UUID]struct{})

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}

		if _, exists := seen[item.ProductID]; exists {
			return ErrDuplicateProduct
		}

		seen[item.ProductID] = struct{}{}
	}

	// load products

	products := make(map[uuid.UUID]*productDomain.Product)

	for _, item := range req.Items {
		product, err := s.store.Products().GetProductByID(ctx, item.ProductID)
		if err != nil {
			return err
		}

		products[item.ProductID] = product
	}

	// verify inventory

	inventories := make(map[uuid.UUID]*inventoryDomain.Inventory)

	for _, item := range req.Items {
		inventory, err := s.store.Inventory().GetInventoryByProductID(ctx, item.ProductID)
		if err != nil {
			return nil
		}

		if inventory.Quantity < item.Quantity {
			return ErrInsufficientInventory
		}

		inventories[item.ProductID] = inventory
	}

	// calculate total

	total := money.Zero()

	for _, item := range req.Items {
		product := products[item.ProductID]

		subTotal := product.Price.Multiply(item.Quantity)

		total = total.Add(subTotal)
	}

	//build order

	now := time.Now()

	order := &domain.Order{
		ID:         uuid.New(),
		UserID:     userID,
		Status:     domain.StatusPending,
		TotalPrice: total,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	items := make([]*domain.OrderItem, 0, len(req.Items))

	for _, reqItem := range req.Items {
		product := products[reqItem.ProductID]

		item := &domain.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: product.ID,

			Quantity: reqItem.Quantity,
			Price:    product.Price,

			CreatedAt: now,
			UpdatedAt: now,
		}

		items = append(items, item)
	}

	// begin transaction
	tx, err := s.store.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// create order
	err = tx.Orders().CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	//create order items
	err = tx.Orders().CreateOrderItems(ctx, items)
	if err != nil {
		return err
	}

	// decrease the inventory
	for _, item := range items {
		err := tx.Orders().DecreaseInventory(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// commit
	return tx.Commit(ctx)

}

func (s *Service) GetOrdersByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]*domain.Order, error) {
	return s.store.Orders().GetOrdersByUserID(ctx, userID)
}

func (s *Service) GetOrderByID(
	ctx context.Context,
	userID uuid.UUID,
	orderID uuid.UUID,
) (*domain.Order, error) {
	return s.store.Orders().GetOrderByID(ctx, userID, orderID)
}
