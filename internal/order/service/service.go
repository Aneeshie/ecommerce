package service

import (
	"context"
	"time"

	"github.com/Aneeshie/ecommerce/internal/common/money"
	inventoryDomain "github.com/Aneeshie/ecommerce/internal/inventory/domain"
	"github.com/Aneeshie/ecommerce/internal/order"
	"github.com/Aneeshie/ecommerce/internal/order/domain"
	"github.com/Aneeshie/ecommerce/internal/order/dto"
	productDomain "github.com/Aneeshie/ecommerce/internal/product/domain"
	"github.com/Aneeshie/ecommerce/internal/store"
	"github.com/google/uuid"
)

var ()

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	CreateOrderItems(ctx context.Context, orderItems []*domain.OrderItem) error
	DecreaseInventory(ctx context.Context, productID uuid.UUID, quantity int) error
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*domain.Order, error)
}

type ProductRepository interface {
	GetProductByID(ctx context.Context, productId uuid.UUID) (*productDomain.Product, error)
}

type InventoryRepository interface {
	GetInventoryByProductID(ctx context.Context, productID uuid.UUID) (*inventoryDomain.Inventory, error)
}

type TxStore interface {
	Orders() OrderRepository
	Products() ProductRepository
	Inventory() InventoryRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Store interface {
	Orders() OrderRepository
	Products() ProductRepository
	Inventory() InventoryRepository
	Begin(ctx context.Context) (TxStore, error)
}

type storeWrapper struct {
	*store.Store
}

func (w *storeWrapper) Orders() OrderRepository {
	return w.Store.Orders()
}

func (w *storeWrapper) Products() ProductRepository {
	return w.Store.Products()
}

func (w *storeWrapper) Inventory() InventoryRepository {
	return w.Store.Inventory()
}

func (w *storeWrapper) Begin(ctx context.Context) (TxStore, error) {
	tx, err := w.Store.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx}, nil
}

type txWrapper struct {
	*store.TxStore
}

func (t *txWrapper) Orders() OrderRepository {
	return t.TxStore.Orders()
}

func (t *txWrapper) Products() ProductRepository {
	return t.TxStore.Products()
}

func (t *txWrapper) Inventory() InventoryRepository {
	return t.TxStore.Inventory()
}

type Service struct {
	store Store
}

func NewService(s *store.Store) *Service {
	return &Service{
		store: &storeWrapper{s},
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, req *dto.CreateOrderRequest) error {
	// validate the request
	if len(req.Items) == 0 {
		return order.ErrEmptyOrder
	}
	seen := make(map[uuid.UUID]struct{})

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return order.ErrInvalidQuantity
		}

		if _, exists := seen[item.ProductID]; exists {
			return order.ErrDuplicateProduct
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
			return order.ErrInsufficientInventory
		}

		inventories[item.ProductID] = inventory
	}

	// calculate total

	total := money.Zero()

	for _, item := range req.Items {
		product := products[item.ProductID]

		var err error
		subTotal, err := product.Price.Multiply(item.Quantity)
		if err != nil {
			return err
		}
		total, err = total.Add(subTotal)
		if err != nil {
			return err
		}
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
