# Users

Guest

Customer

Admin

---

# Functional Requirements

Authentication

Catalog

Shopping

Checkout

Orders

Admin

---

# Non Functional Requirements

Performance

Security
- Authentication via JSON Web Tokens (JWT)
- Argon2 / Bcrypt strong password hashing
- Role-Based Access Control (RBAC) middleware

Reliability
- Database transactional rollbacks on partial failures
- Anti-corruption layer mapping domain errors to HTTP codes (httpx)

Scalability

Observability

Integration Testing
- Strict avoidance of Mock databases
- Ephemeral PostgreSQL instances via Testcontainers
- Transaction-rollback isolation for every test

Continuous Integration (CI/CD)
- Fully deterministic environments via Nix Flakes
- Automated GitHub Actions pipeline (`go vet`, `go build`, `go test`)

Modular Architecture
- Domain-Driven Design (DDD) module boundaries
- Store Wrapper pattern for dependency injection
