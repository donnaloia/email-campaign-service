package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:     getEnvOrDefault("DB_HOST", "db"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "sendpulse"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func Connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Initialize database with migrations
	if err := initializeDatabase(db); err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	return db, nil
}

func initializeDatabase(db *sql.DB) error {
	fmt.Println("initializing database")
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Get all .sql files
	migrationPath := "/app/migrations/*.sql"
	files, err := filepath.Glob(migrationPath)
	if err != nil {
		return fmt.Errorf("error finding migration files: %w", err)
	}
	for _, file := range files {
		fmt.Printf("- %s\n", file)
	}
	sort.Strings(files)

	for _, file := range files {
		// Skip .down.sql files
		if strings.Contains(file, "down") {
			continue
		}

		// Check if migration has already been executed
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)", file).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking migration status: %w", err)
		}

		if exists {
			fmt.Printf("Skipping already executed migration: %s\n", file)
			continue
		}

		// Read and execute the migration
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", file, err)
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		// Execute the migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing %s: %w", file, err)
		}

		// Record the migration
		_, err = tx.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", file)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", file, err)
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing transaction: %w", err)
		}

		fmt.Printf("Executed migration: %s\n", file)
	}

	return nil
}
