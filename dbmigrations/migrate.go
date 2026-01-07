package dbmigrations

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB, command string) error {
	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set migration directory
	migrationDir := "./dbmigrations/migrations/mysql"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return fmt.Errorf("migration directory does not exist: %s", migrationDir)
	}

	// Set provider to mysql
	goose.SetDialect("mysql")

	// Run goose command
	switch command {
	case "up":
		if err := goose.Up(sqlDB, migrationDir); err != nil {
			return fmt.Errorf("failed to run migrations up: %w", err)
		}
	case "down":
		if err := goose.Down(sqlDB, migrationDir); err != nil {
			return fmt.Errorf("failed to run migrations down: %w", err)
		}
	case "status":
		if err := goose.Status(sqlDB, migrationDir); err != nil {
			return fmt.Errorf("failed to get migration status: %w", err)
		}
	default:
		return fmt.Errorf("unknown migration command: %s", command)
	}

	return nil
}
