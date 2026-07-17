# Orders Design

## Goal

The Orders module is responsible for recording purchases made by users.

An order represents a single checkout operation. During checkout, a user may purchase one or more products in a single order.

The module is responsible for:

- Validating requested products
- Verifying inventory availability
- Recording the purchase
- Preserving historical pricing
- Updating inventory atomically

---

# Responsibilities

An order should:

- Belong to exactly one user.
- Contain one or more products.
- Preserve the purchase price of every item.
- Record the total amount paid.
- Track the current order status.
- Be created atomically alongside inventory updates.

---

# Why Two Tables?

At first glance, it may seem reasonable to store a `product_id` inside the `orders` table.

```text
orders

id
user_id
product_id
quantity
```

However, this design only supports a single product per order.

Real-world orders often contain multiple products.

Example:

Order #123

- Aula F75 ×1
- Logitech G304 ×2
- Mouse Pad ×1

Because an order can contain many products, the database models them separately.

An order represents the purchase itself.

Order items represent the individual products inside that purchase.

---

# Schema

## Orders

Stores information about the order as a whole.

```sql
orders

id
user_id
status
total_price

created_at
updated_at
```

### Fields

| Field | Description |
|--------|-------------|
| id | Unique order identifier |
| user_id | Customer who placed the order |
| status | Current order status |
| total_price | Total value of the order |
| created_at | Creation timestamp |
| updated_at | Last update timestamp |

Notice that the table does **not** contain a `product_id`.

Products belong to the order through the `order_items` table.

---

## Order Items

Stores every product purchased within an order.

```sql
order_items

id
order_id
product_id

quantity
price

created_at
updated_at
```

Each row represents one purchased product.

Example:

| order_id | product | quantity | price |
|----------|---------|---------:|------:|
|123|Aula F75|1|6999|
|123|Logitech G304|2|3999|

One order can therefore contain any number of products.

---

# Why Store the Price?

Product prices change over time.

Example:

Today:

Aula F75 → ₹6999

Next Month:

Aula F75 → ₹7999

If an order referenced the current product price, historical orders would become incorrect.

Instead, each order item stores the purchase price at the time the order was placed.

This preserves historical accuracy.

---

# Why Store the Total?

The total price is stored in the `orders` table instead of being calculated on every request.

Reasons:

- Orders are read frequently.
- The total never changes after creation.
- Reading a stored value is more efficient than recalculating it every time.

The total is calculated once during checkout and then persisted.

---

# Relationships

```text
Users
  │
  │ 1 → many
  ▼
Orders
  │
  │ 1 → many
  ▼
Order Items
  │
  │ many → 1
  ▼
Products
```

A user may place many orders.

Each order contains one or more order items.

Each order item references exactly one product.

---

# Order Status

For Version 1, the following statuses are supported:

```text
Pending
Paid
Cancelled
```

Additional statuses such as Processing, Shipped, Delivered, and Refunded can be introduced in future iterations.

---

# API

## Create Order

```http
POST /orders
```

Request

```json
{
  "items": [
    {
      "productId": "uuid",
      "quantity": 2
    },
    {
      "productId": "uuid",
      "quantity": 1
    }
  ]
}
```

The client only specifies:

- Product
- Quantity

The server determines:

- User (from authentication)
- Product prices
- Order total
- Order status

The client is never trusted to provide monetary values.

---

# Business Rules

- User must be authenticated.
- Order must contain at least one item.
- Every product must exist.
- Every product must be active.
- Quantity must be greater than zero.
- Requested quantity must not exceed available inventory.
- Prices are copied from the product at checkout.
- Order total is calculated by the server.
- Order creation must be transactional.

---

# Transaction Flow

Creating an order is an atomic operation.

```text
Begin Transaction

↓

Validate User

↓

Validate Products

↓

Validate Inventory

↓

Calculate Total

↓

Create Order

↓

Create Order Items

↓

Reduce Inventory

↓

Commit
```

If any step fails, the transaction is rolled back.

This guarantees that:

- No partial orders exist.
- Inventory remains consistent.
- Order items always belong to a valid order.

---

# Future Enhancements

The initial implementation focuses on the core order lifecycle.

Future improvements may include:

- Shipping addresses
- Coupons and discounts
- Taxes
- Payment integration
- Order history
- Returns and refunds
- Shipment tracking
- Inventory reservations
- Event-driven order processing
