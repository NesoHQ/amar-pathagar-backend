package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Connection wraps sql.DB
type Connection struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(ctx context.Context, connectionString string) (*Connection, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection with context
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.DB.Close()
}
