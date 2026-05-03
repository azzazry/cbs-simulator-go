package database

import (
	"cbs-simulator/config"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

// InitDB initializes the PostgreSQL connection and runs migrations
func InitDB() error {
	var err error

	// Open database connection via pgx stdlib driver
	DB, err = sql.Open("pgx", config.AppConfig.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Connection pool tuning
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("PostgreSQL connection established")

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
		"./database/migrations/005_add_security_tables.sql",
		"./database/migrations/006_core_banking.sql",
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
	seederFiles := []string{
		"./database/seeders/001_sample_data.sql",
		"./database/seeders/002_sample_banks.sql",
		"./database/seeders/003_sample_transfer_fees.sql",
		"./database/seeders/004_security_seed.sql",
		"./database/seeders/005_core_banking_seed.sql",
	}

	for _, seederFile := range seederFiles {
		sqlBytes, err := os.ReadFile(seederFile)
		if err != nil {
			log.Printf("Warning: seeder file not found: %s", seederFile)
			continue
		}

		if _, err := DB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute seeder %s: %v", seederFile, err)
		}
	}

	log.Println("Seeders completed successfully")
	return nil
}

// isEmpty checks if the database is empty (no customers yet)
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
