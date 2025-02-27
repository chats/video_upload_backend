package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"cams.dev/video_upload_backend/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get applied migrations
	rows, err := db.Query("SELECT name FROM migrations ORDER BY id")
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Failed to scan migration name: %v", err)
		}
		appliedMigrations[name] = true
	}

	// Get migration files
	migrationFiles, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatalf("Failed to get migration files: %v", err)
	}

	// Sort migration files by name
	sort.Strings(migrationFiles)

	// Apply migrations
	for _, file := range migrationFiles {
		fileName := filepath.Base(file)

		// Skip if already applied
		if appliedMigrations[fileName] {
			log.Printf("Skipping already applied migration: %s", fileName)
			continue
		}

		// Read migration file
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", fileName, err)
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to start transaction: %v", err)
		}

		// Execute migration
		log.Printf("Applying migration: %s", fileName)

		// Split by semicolon to handle multiple statements
		statements := strings.Split(string(content), ";")
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			_, err = tx.Exec(statement)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to execute migration %s: %v", fileName, err)
			}
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (name, applied_at) VALUES ($1, $2)", fileName, time.Now())
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to record migration %s: %v", fileName, err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		log.Printf("Migration applied: %s", fileName)
	}

	log.Println("Migrations completed successfully!")
}
