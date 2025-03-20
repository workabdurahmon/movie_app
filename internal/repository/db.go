package repository

import (
	"fmt"
	"movie_app/internal/config"
	"movie_app/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Apply auto-migrations for schema updates
	if err := db.AutoMigrate(&domain.Movie{}, &domain.User{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	// Get the underlying SQL database connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pooling for scalability
	sqlDB.SetMaxOpenConns(100)                // Maximum number of open connections
	sqlDB.SetMaxIdleConns(20)                 // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum time a connection can be reused

	// Test database connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
