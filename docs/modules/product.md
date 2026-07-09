# Product Module

## Goal

The Product module is responsible for managing the product catalog.

It allows administrators to create and manage products, while allowing customers to browse and view them.

---

# Functional Requirements

## Administrator

- Create a product.
- Update a product.
- Delete a product.

## Customer

- View a product.
- View all products.

---

# Product

A product consists of:

- Name
- Description
- Price

---

# Business Rules

- Product name cannot be empty.
- Product description cannot be empty.
- Product price cannot be negative.
- Only administrators can create, update or delete products.

---

# Out of Scope (V1)

These features will be implemented in future modules.

- Product Categories
- Product Images
- Product Inventory / Stock
- Product Variants
- SKU
- Product Reviews
- Product Ratings
- Search
- Pagination
- Filtering
- Discounts

---

# API

## Admin

POST   /products

PUT    /products/{id}

DELETE /products/{id}

## Customer

GET    /products

GET    /products/{id}

---

# Database

## products

- id
- name
- description
- price
- status
- created_at
- updated_at

---

# Business Rules

- Product name cannot be empty.
- Product description cannot be empty.
- Product price cannot be negative.
- Only administrators can create, update or delete products.
- Products with existing orders cannot be permanently deleted and should instead be archived.

---

# Product Status

Products can have one of the following states:

- ACTIVE
- ARCHIVED

Archived products are retained for historical purposes but are not available for customers to browse or purchase.

---
# Notes
This is intentionally a minimal Product module.

The focus of V1 is to establish the Product domain and its CRUD operations.

Inventory, Categories, Images, SKU, Search and other ecommerce features will be introduced as separate modules to keep the design simple and allow the system to evolve incrementally.
