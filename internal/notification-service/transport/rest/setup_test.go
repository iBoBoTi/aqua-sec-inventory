package rest_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/domain"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createPostgresContainer(t *testing.T, dbName, dbUser, dbPassword string, logger *slog.Logger) (string, string) {
	t.Helper()
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
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

	return host, port.Port()
}

func setUpTestDB(t *testing.T, dbName, dbUser, dbPassword string) (*sql.DB, error) {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	host, port := createPostgresContainer(t, dbName, dbUser, dbPassword, logger)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS resources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customer_resource (
    id SERIAL PRIMARY KEY, -- Optional auto-increment ID
    customer_id INT NOT NULL,
    resource_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(), -- Corrected default value
    UNIQUE (customer_id, resource_id), -- Prevent duplicate relationships
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`)
	assert.NoError(t, err)

	return db, nil
}

func seedNotification(t *testing.T, db *sql.DB) domain.Notification {

	var notification domain.Notification
	query := `
		INSERT INTO notifications (user_id, message) 
		VALUES ($1, $2) 
		RETURNING id, user_id, message;
	`

	err := db.QueryRow(query, 1, "added a cloud resource").Scan(&notification.ID, &notification.UserID, &notification.Message)
	assert.NoError(t, err)

	return notification
}
