package database

import (
	"cbs-simulator/config"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the database connection and creates tables
func InitDB() error {
	var err error

	// Ensure database directory exists
	dbDir := "./database"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %v", err)
		}
	}

	// Open database connection
	DB, err = sql.Open("sqlite", config.AppConfig.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Check if database is empty, then seed
	if isEmpty() {
		log.Println("Database is empty, running seeders...")
		if err := runSeeders(); err != nil {
			return fmt.Errorf("failed to run seeders: %v", err)
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

// runMigrations executes all SQL migration files
func runMigrations() error {
	migrationFiles := []string{
		"./database/migrations/001_init_schema.sql",
		"./database/migrations/002_add_notifications.sql",
		"./database/migrations/003_add_banks.sql",
		"./database/migrations/004_add_transfer_fees.sql",
	}

	for _, migrationFile := range migrationFiles {
		sqlBytes, err := os.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", migrationFile, err)
		}

		if _, err := DB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", migrationFile, err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}

// runSeeders executes all SQL seeder files
func runSeeders() error {
	seederFile := "./database/seeders/001_sample_data.sql"

	sqlBytes, err := os.ReadFile(seederFile)
	if err != nil {
		return fmt.Errorf("failed to read seeder file: %v", err)
	}

	if _, err := DB.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to execute seeder: %v", err)
	}

	log.Println("Seeders completed successfully")
	return nil
}

// isEmpty checks if the database is empty
func isEmpty() bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return true
	}
	return count == 0
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
