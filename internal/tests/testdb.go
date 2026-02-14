package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// TestDB wraps the test database connection and provides helper methods
type TestDB struct {
	Container        *PostgresContainer
	Pool             *pgxpool.Pool
	ConnectionString string
}

// SetupTestDB creates a new test database with migrations applied
func SetupTestDB(t *testing.T, ctx context.Context) (*TestDB, func()) {
	t.Helper()

	// Create postgres container
	pgContainer, err := CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Fatalf("failed to create postgres container: %v", err)
	}

	// Create connection pool
	poolConfig, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		t.Fatalf("failed to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to create connection pool: %v", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	testDB := &TestDB{
		Container:        pgContainer,
		Pool:             pool,
		ConnectionString: pgContainer.ConnectionString,
	}

	// Cleanup function
	cleanup := func() {
		pool.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}

	return testDB, cleanup
}

// Truncate removes all data from the specified tables
func (db *TestDB) Truncate(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if _, err := db.Pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}

// TruncateAll removes all data from all tables
func (db *TestDB) TruncateAll(ctx context.Context) error {
	return db.Truncate(ctx, "thumbnails", "files", "users")
}
