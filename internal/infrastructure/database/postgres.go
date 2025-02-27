package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresDB provides a connection to a PostgreSQL database
type PostgresDB struct {
	db *sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(config Config) (*PostgresDB, error) {
	// Format connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

// DB returns the underlying database connection
func (p *PostgresDB) DB() *sql.DB {
	return p.db
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.db.Close()
}
