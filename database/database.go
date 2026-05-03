package database

import (
	"cbs-simulator/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() error {
	var err error

	DB, err = sql.Open("pgx", config.AppConfig.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	if isEmpty() {
		if err := runSeeders(); err != nil {
			return fmt.Errorf("failed to run seeders: %v", err)
		}
	}

	return nil
}

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

	return nil
}

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
			continue // skip kalau file tidak ada
		}

		if _, err := DB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute seeder %s: %v", seederFile, err)
		}
	}

	return nil
}

func isEmpty() bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return true
	}
	return count == 0
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
