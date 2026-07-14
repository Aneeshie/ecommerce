# ADR 0003 — Introduce Store Pattern

## Status

Accepted

---

## Context

As the application grows, business operations may require multiple repositories to participate in the same database transaction.

For example:

- Create Product
- Create Inventory

Both operations must either succeed together or fail together.

The existing architecture injects repositories directly into services.

```text
ProductService
        │
        ▼
ProductRepository
```

This design makes it difficult to coordinate multiple repositories within a single transaction.

---

## Decision

Introduce a Store that owns the database connection pool.

The Store is responsible for:

- Creating repositories
- Beginning transactions

A transaction produces a Transaction Store (`TxStore`), which creates repositories backed by the same database transaction.

---

## Architecture

```text
                Service
                   │
                   ▼
                Store
          ┌────────┴────────┐
          ▼                 ▼
      Database Pool      Transaction
          │                 │
          ▼                 ▼
 Product Repository   Product Repository
 Inventory Repository Inventory Repository
```

---

## Benefits

- Services own business workflows.
- Repositories remain responsible only for persistence.
- Multiple repositories can participate in the same transaction.
- Future modules (Orders, Payments, Inventory) reuse the same transaction infrastructure.

---

## Alternatives Considered

### Inject the database pool into every service

This allows services to begin transactions directly.

However, services must also construct repositories manually, leading to duplicated infrastructure code.

---

### Let repositories own transactions

Rejected.

Repositories should not coordinate business workflows involving multiple modules.

---

## Consequences

Services now depend on the Store rather than individual repositories.

Transactions become an infrastructure concern managed by the Store.

Repositories remain unaware of whether they are using a connection pool or an active transaction.
