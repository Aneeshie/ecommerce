# Ecommerce Backend

A production-inspired ecommerce backend built in Go with a focus on clean architecture, domain-driven design principles, and modular development.

The project is being developed module-by-module, following industry conventions such as layered architecture, feature-based modules, JWT authentication, role-based authorization, and PostgreSQL.

---

## Tech Stack

- Go
- PostgreSQL
- pgx
- Chi Router
- JWT Authentication
- Docker
- UUID

---

## Architecture

The project follows a layered architecture.

```
HTTP
│
├── Handler
│
├── Service
│
├── Repository
│
└── PostgreSQL
```

Each module is self-contained and consists of:

```
internal/
└── <module>/
    ├── domain/
    ├── dto/
    ├── repository/
    ├── service/
    └── handler/
```

---

## Project Structure

```
cmd/
└── api/

internal/
├── common/
│   └── money/
├── identity/
├── middleware/
└── product/

migrations/
```

---

## Features

### Identity

- User Registration
- User Login
- Refresh Tokens
- JWT Authentication
- Current User Endpoint (`/auth/me`)
- Role-Based Authorization Middleware

### Product

- Create Product
- List Products
- Get Product
- Update Product
- Archive Product (Soft Delete)

---

## API

### Authentication

| Method | Endpoint |
|---------|----------|
| POST | `/api/v1/auth/register` |
| POST | `/api/v1/auth/login` |
| POST | `/api/v1/auth/refresh` |
| GET | `/api/v1/auth/me` |

### Products

| Method | Endpoint | Access |
|---------|----------|--------|
| POST | `/api/v1/products` | Admin |
| GET | `/api/v1/products` | Public |
| GET | `/api/v1/products/{id}` | Public |
| PUT | `/api/v1/products/{id}` | Admin |
| DELETE | `/api/v1/products/{id}` | Admin |

---

## Business Rules

### Users

- Email addresses must be unique.
- Passwords are securely hashed.
- JWT access tokens are used for authentication.
- Refresh tokens can be exchanged for new access tokens.

### Products

- Product prices are stored using a dedicated `Money` value object.
- Product prices cannot be negative.
- Products are never permanently deleted.
- Archived products are excluded from the public catalogue.
- Only administrators can create, update, or archive products.

---

## Value Objects

### Money

Instead of using floating point numbers, prices are stored as integer values representing the smallest currency unit (e.g. paise).

Example:

```
₹2,499.99
↓

249999
```

This avoids floating point precision issues.

---

## Running

```bash
git clone <repository>

cd ecommerce

go run cmd/api/main.go
```

---

## Roadmap

- [x] Identity Module
- [x] Product Module
- [ ] Inventory Module
- [ ] Shopping Cart
- [ ] Orders
- [ ] Payments
- [ ] Admin Dashboard
- [ ] Frontend

---

## Goals

This project is intended as a long-term learning project to explore production backend development concepts, including:

- Clean Architecture
- Domain-Driven Design
- Repository Pattern
- JWT Authentication
- Role-Based Authorization
- Database Design
- Transactions
- Concurrency
- Testing
- Distributed Systems Concepts

---

## License

MIT
