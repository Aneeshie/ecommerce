package testdb

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
)

func initPool(ctx context.Context) error {
	var err error
	poolOnce.Do(func() {
		// Use a local database for testing since Docker is not available.
		// Fallback to a default local postgres instance if TEST_DATABASE_URL is not provided.
		dbURL := os.Getenv("TEST_DATABASE_URL")
		if dbURL == "" {
			// Assuming a local postgres instance with a test database
			dbURL = "postgres://postgres:postgres@localhost:5432/ecommerce_test?sslmode=disable"
		}

		// Run migrations
		err = runMigrations(dbURL)
		if err != nil {
			return
		}

		config, errParse := pgxpool.ParseConfig(dbURL)
		if errParse != nil {
			err = errParse
			return
		}

		pool, err = pgxpool.NewWithConfig(ctx, config)
	})
	return err
}

func runMigrations(dbURL string) error {
	_, b, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(b), "../../..", "migrations")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// SetupTestDB initializes the test database (if not already done)
// and returns a transaction that will automatically be rolled back
// when the test completes. This ensures test isolation.
func SetupTestDB(t *testing.T) pgx.Tx {
	t.Helper()
	ctx := context.Background()

	err := initPool(ctx)
	if err != nil {
		t.Fatalf("failed to init test database pool: %v", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})

	return tx
}
