# Inventory Module

The Inventory module is responsible for tracking the quantity of each product that is currently available for sale.

For Version 1, inventory is tracked globally. Warehouses, locations, and stock distribution are intentionally out of scope and may be introduced in future versions.

---

# Purpose

Inventory determines whether a product can be purchased.

It is responsible for:

- Tracking stock quantity
- Updating stock levels
- Preventing negative inventory
- Providing current stock information

---

# Data Model

An inventory record consists of:

- Product ID
- Quantity
- Created At
- Updated At

Each inventory record belongs to exactly one product.

---

# Business Actions

The Inventory module supports the following operations:

- Create inventory
- Get inventory
- Update inventory (restock or adjust quantity)

Inventory records are never deleted.

---

# Business Rules

- One product has exactly one inventory record.
- Inventory cannot exist without a product.
- Quantity cannot be negative.
- Only administrators can create or update inventory.
- Customers may view inventory availability.
- Inventory is tracked globally for Version 1.
- Warehouses and location-aware inventory are future enhancements.

---

# API

| Method | Endpoint | Access |
|---------|----------|--------|
| POST | `/api/v1/products/{id}/inventory` | Admin |
| GET | `/api/v1/products/{id}/inventory` | Public |
| PUT | `/api/v1/products/{id}/inventory` | Admin |

---

# Future Enhancements

The following features are intentionally excluded from Version 1:

- Multi-warehouse inventory
- Stock reservations
- Automatic stock deduction during checkout
- Inventory transfers between warehouses
- Inventory history and audit logs
- Low-stock notifications
- Reserved vs Available inventory

---

# Notes

Inventory represents the quantity currently available for sale.

Version 1 stores a single global quantity per product. Future versions may associate inventory with warehouses or fulfillment centers while preserving the same business concept.

# Dependencies

## Depends On

- Product Module

## Used By

- Cart Module
- Order Module