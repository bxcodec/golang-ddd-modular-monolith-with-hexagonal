package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
	ConnStr   string
}

func SetupPostgres(t *testing.T) *PostgresContainer {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	return &PostgresContainer{
		Container: postgresContainer,
		DB:        db,
		ConnStr:   connStr,
	}
}

func (pc *PostgresContainer) RunMigrations(t *testing.T, migrationsPath string) {
	absPath, err := filepath.Abs(migrationsPath)
	require.NoError(t, err)

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		pc.ConnStr,
	)
	require.NoError(t, err)

	err = migrator.Up()
	require.NoError(t, err)

	sourceErr, dbErr := migrator.Close()
	require.NoError(t, sourceErr)
	require.NoError(t, dbErr)
}

func (pc *PostgresContainer) TruncateTables(t *testing.T, tables ...string) {
	for _, table := range tables {
		_, err := pc.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err)
	}
}

func (pc *PostgresContainer) Teardown(t *testing.T) {
	ctx := context.Background()

	if pc.DB != nil {
		err := pc.DB.Close()
		require.NoError(t, err)
	}

	if pc.Container != nil {
		err := pc.Container.Terminate(ctx)
		require.NoError(t, err)
	}
}
