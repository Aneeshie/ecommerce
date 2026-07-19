<div align="center">
  <h1>🛒 Aneeshie Ecommerce Backend</h1>
  <p>A production-grade, domain-driven ecommerce API built in Go.</p>

  <!-- Badges -->
  <p>
    <a href="https://github.com/Aneeshie/ecommerce/actions/workflows/ci.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/Aneeshie/ecommerce/ci.yml?branch=main&style=flat-square&logo=github" alt="Build Status" />
    </a>
    <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" alt="Go Version" />
    <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=flat-square&logo=postgresql" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/Docker-Testcontainers-2496ED?style=flat-square&logo=docker" alt="Docker Testcontainers" />
    <img src="https://img.shields.io/badge/Architecture-DDD-orange?style=flat-square" alt="DDD Architecture" />
  </p>
</div>

---

## 📖 Overview

This project is a high-performance ecommerce backend designed as a rigorous exploration of production-ready development concepts. It avoids tightly-coupled monolithic designs in favor of **Modular Monolith** and **Domain-Driven Design (DDD)** principles.

Every feature (Identity, Products, Orders, Inventory) lives in its own isolated module with strict boundaries, layered architecture, and comprehensive integration testing.

---

## ⚡ Key Features

- **Domain-Driven Modules**: Self-contained modules (`identity`, `product`, `order`) with strict layered boundaries.
- **Robust Integration Testing**: Ephemeral, throwaway database instances generated per-test using **Testcontainers**, guaranteeing 100% test isolation via SQL transactions.
- **Secure Authentication**: Stateless, JWT-based authentication with Access and Refresh tokens, backed by Argon2/Bcrypt password hashing and Role-Based Access Control (RBAC).
- **Financial Precision**: Custom `Money` value objects storing currency as exact integer cents/paise to mathematically eliminate floating-point precision errors.
- **Reproducible Environments**: Powered by **Nix** (`flake.nix`), guaranteeing that the development environment on your Mac matches the GitHub Actions CI server down to the exact byte.

---

## 🛠 Tech Stack

| Category | Technology | Purpose |
|----------|------------|---------|
| **Core** | [Go (Golang)](https://go.dev/) | High-performance, statically typed backend logic |
| **Routing** | [go-chi](https://github.com/go-chi/chi) | Lightweight, idiomatic HTTP routing |
| **Database** | [PostgreSQL](https://www.postgresql.org/) | Relational data persistence |
| **DB Driver** | [pgx](https://github.com/jackc/pgx) | High-performance Go driver for PostgreSQL |
| **Testing** | [Testcontainers-Go](https://golang.testcontainers.org/) | Docker-backed integration testing |
| **Migrations**| [golang-migrate](https://github.com/golang-migrate/migrate) | Programmatic schema migrations |
| **CI/CD** | [GitHub Actions](https://github.com/features/actions) | Automated linting, building, and testing |
| **Tooling** | [Nix](https://nixos.org/) | Deterministic, reproducible development shells |

---

## 🏗 Architecture

The application strictly adheres to a **Layered Architecture** within isolated modules:

```text
HTTP Request
     │
     ▼
┌──────────────┐    Parses JSON, validates input, manages HTTP status codes.
│   Handler    │    Knows nothing about business logic or databases.
└──────┬───────┘
       ▼
┌──────────────┐    The brain of the operation. Orchestrates business rules,
│   Service    │    handles domain logic, and coordinates multiple repositories.
└──────┬───────┘
       ▼
┌──────────────┐    Handles all raw SQL queries, schema mapping, and pgx
│  Repository  │    transactions. Knows nothing about HTTP or business context.
└──────┬───────┘
       ▼
  PostgreSQL 
```

---

## 🚀 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) (or [OrbStack](https://orbstack.dev/) for macOS)
- [Nix](https://nixos.org/download) (Optional, but highly recommended for the dev shell)
- Go 1.22+

### 1. Clone & Setup
```bash
git clone git@github.com:Aneeshie/ecommerce.git
cd ecommerce
```

### 2. Enter the Development Shell
If you are using Nix, drop into the deterministic shell which will automatically install Go, Make, and all dependencies:
```bash
nix develop
```

### 3. Start the Local Database
Spin up the development PostgreSQL database in the background:
```bash
docker-compose up -d
```

### 4. Run the API
```bash
go run cmd/api/main.go
```

---

## 🧪 Testing Strategy

We do not use Mock databases. Instead, we test against **reality**. 

The test suite leverages `testcontainers-go` to dynamically spin up an isolated PostgreSQL Docker container exclusively for the duration of the tests. 
- It automatically applies `golang-migrate` schemas.
- Every individual repository test is wrapped in an SQL `Tx` (Transaction).
- When the test asserts its results, the transaction is immediately **rolled back**, guaranteeing a pristine database for the next test.

**Run the tests:**
```bash
nix develop --command go test -v ./...
```
*(Note: If you use OrbStack on macOS, the test suite is configured to automatically bridge the Unix socket connection!)*

---

## 🔌 API Reference

### Authentication (`/api/v1/auth`)
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| `POST` | `/register` | Register a new user account | Public |
| `POST` | `/login` | Authenticate and retrieve JWTs | Public |
| `POST` | `/refresh` | Exchange a refresh token for an access token | Public |
| `GET`  | `/me` | Retrieve the authenticated user's profile | **Auth** |

### Products (`/api/v1/products`)
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| `GET`  | `/` | List all active products | Public |
| `GET`  | `/{id}` | Retrieve a specific product | Public |
| `POST` | `/` | Create a new product | **Admin** |
| `PUT`  | `/{id}` | Update product details | **Admin** |
| `DELETE`| `/{id}` | Archive a product (Soft delete) | **Admin** |

### Orders (`/api/v1/orders`)
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| `POST` | `/` | Place a new order | **Auth** |
| `GET`  | `/` | List all orders for the current user | **Auth** |
| `GET`  | `/{id}` | Get specific order details | **Auth** |

---

## 🗺 Roadmap

- [x] **Identity Module** (Users, JWT, Roles, Passwords)
- [x] **Product Module** (Catalog, Soft Deletes, Admin Guards)
- [x] **Orders Module** (Placing Orders, Line Items)
- [x] **Testing Infrastructure** (Testcontainers, TX Rollbacks)
- [x] **CI/CD** (GitHub Actions with Nix integration)
- [ ] **Inventory Module** (Stock tracking, reservations)
- [ ] **Shopping Cart** (Redis-backed session carts)
- [ ] **Payments** (Stripe integration)

---

## 📄 License
This project is licensed under the MIT License.
