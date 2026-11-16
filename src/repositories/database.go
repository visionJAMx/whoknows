package repositories

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Open connects to PostgreSQL using a connection string
func Open(connectionString string) (*sql.DB, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("database connection string must not be empty")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// PostgreSQL connection settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	return db, nil
}

// Initialize is a no-op since Alembic handles schema
func Initialize(ctx context.Context, db *sql.DB) error {
	// Schema is managed by Alembic migrations
	return nil
}
