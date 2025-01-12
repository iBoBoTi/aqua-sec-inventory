package rest_test

import (
	"context"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TODO: consider moving this into a different package of its own. maybe a package called testing
// or similar to put reusable test functions
func CreatePostgresContainer(t *testing.T, dbName, dbUser, dbPassword string, logger *slog.Logger) (string, uint) {
	t.Helper()
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		// postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err)
	})

	host, err := postgresContainer.Host(context.Background())
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(context.Background(), "5432/tcp")
	require.NoError(t, err)

	dbport, err := strconv.Atoi(port.Port())
	require.NoError(t, err)

	logger.Info("container info", "host", host, "port", dbport)

	return host, uint(dbport)
}
