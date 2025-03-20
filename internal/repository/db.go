package repository

import (
	"movie_app/internal/config"
	"movie_app/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&domain.Movie{}, &domain.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
