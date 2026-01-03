package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// Connect establishes a connection to PostgreSQL using DATABASE_URL
func Connect(databaseURL string) (*pgx.Conn, error) {
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set. Set it or use --database-url flag")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return conn, nil
}

// CheckExtension verifies pg_migrate extension is installed
func CheckExtension(conn *pgx.Conn) error {
	ctx := context.Background()

	var exists bool
	err := conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'pg_migrate')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("cannot check extension: %w", err)
	}

	if !exists {
		return fmt.Errorf("pg_migrate extension not installed. Run: CREATE EXTENSION pg_migrate")
	}

	return nil
}

// GetExtensionVersion returns the pg_migrate extension version
func GetExtensionVersion(conn *pgx.Conn) (string, error) {
	ctx := context.Background()

	var version string
	err := conn.QueryRow(ctx,
		"SELECT extversion FROM pg_extension WHERE extname = 'pg_migrate'").Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}
