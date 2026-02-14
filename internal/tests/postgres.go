package tests

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// migrationsDir returns the path to db/migrations relative to the module root.
// When running tests, the working directory is the package under test, so we
// walk up to find go.mod and then use that root + db/migrations.
func migrationsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			p := filepath.Join(dir, "db", "migrations")
			if _, err := os.Stat(p); err != nil {
				return "", err
			}
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(t *testing.T, ctx context.Context) (*PostgresContainer, error) {
	logger := log.Default()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithLogger(logger),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Run migrations
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	migrationsPath, err := migrationsDir()
	if err != nil {
		return nil, err
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
