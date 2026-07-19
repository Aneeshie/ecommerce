# Integration Testing Infrastructure

## Overview

The testing strategy in this project fundamentally avoids **Mocks** at the repository layer. Instead, it relies on true **Integration Testing** against an ephemeral PostgreSQL database.

This guarantees that:
- SQL syntax is strictly verified.
- Database constraints (Unique, Foreign Keys, Checks, Enums) are actually enforced.
- Complex queries and joins behave exactly as they will in production.

---

## Testcontainers

Instead of requiring developers to manually configure a local PostgreSQL instance, the project integrates **[testcontainers-go](https://golang.testcontainers.org/)**.

### How it Works

1. **Docker Engine**: The Go test suite acts as a Docker client and communicates with your local Docker Daemon.
2. **Ephemeral Containers**: It spins up a throwaway `postgres:16-alpine` container specifically for the test run.
3. **Dynamic Ports**: It maps the database to a dynamic, high-numbered port to avoid conflicts with existing databases.
4. **Wait Strategy**: It polls the database logs until the server is fully ready to accept connections before returning control to the test suite.
5. **Garbage Collection (Ryuk)**: A hidden sidecar container (Ryuk) monitors the Go test process. The instant the test process finishes or crashes, Ryuk aggressively deletes the container, volumes, and networks, ensuring no orphan resources remain.

---

## The Transaction Rollback Strategy

Spinning up a Docker container takes roughly ~2 seconds. If we spun up a new container for *every individual test*, the test suite would be incredibly slow.

Instead, we use a **Transaction Rollback Strategy** to guarantee test isolation while maintaining high performance.

### `testdb.SetupTestDB(t)`

Found in `internal/common/testdb/testdb.go`, this helper function is injected into every repository test:

```go
func SetupTestDB(t *testing.T) pgx.Tx {
    // 1. Ensures the Docker container is running (Thread-safe Singleton)
    // 2. Begins a new pgx.Tx transaction
    // 3. Registers a t.Cleanup() hook to rollback the transaction
}
```

### The Test Lifecycle

1. **Suite Start**: The first test triggers the `testcontainers-go` initialization.
2. **Migrations**: `golang-migrate` connects to the dynamic DB port and runs all `.sql` migrations in `/migrations`, establishing the schema.
3. **Test Execution**: A test calls `tx := testdb.SetupTestDB(t)`. The database begins a transaction (`BEGIN;`). The test inserts fixtures, runs queries, and asserts results within that isolated transaction.
4. **Test Cleanup**: When the test finishes, Go executes the `t.Cleanup()` hook, issuing a `ROLLBACK;`. The database discards all changes.
5. **Next Test**: The next test receives a perfectly pristine database state in microseconds.

### Benefits
- 100% data isolation between tests.
- Extremely fast execution times because the schema and container are reused.
- Strict enforcement of foreign keys and relational integrity (e.g., creating an Order in tests requires first inserting a valid User and Product).
