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
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
)

func initPool(ctx context.Context) error {
	var err error
	poolOnce.Do(func() {
		// OrbStack specific fix for Nix environment on macOS
		if runtime.GOOS == "darwin" && os.Getenv("DOCKER_HOST") == "" {
			os.Setenv("DOCKER_HOST", "unix:///Users/nara/.orbstack/run/docker.sock")
		}

		postgresContainer, errStart := tcpostgres.Run(ctx,
			"postgres:16-alpine",
			tcpostgres.WithDatabase("ecommerce_test"),
			tcpostgres.WithUsername("testuser"),
			tcpostgres.WithPassword("testpassword"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		if errStart != nil {
			err = errStart
			return
		}

		connString, errConn := postgresContainer.ConnectionString(ctx, "sslmode=disable")
		if errConn != nil {
			err = errConn
			return
		}

		// Run migrations
		err = runMigrations(connString)
		if err != nil {
			return
		}

		config, errParse := pgxpool.ParseConfig(connString)
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
