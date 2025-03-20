package repository

import (
	"context"
	"errors"
	"fmt"
	"movie_app/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrCreateMovie = errors.New("failed to create movie")
	ErrUpdateMovie = errors.New("failed to update movie")
	ErrDeleteMovie = errors.New("failed to delete movie")
	ErrFetchMovie  = errors.New("failed to fetch movie")
)

type MovieRepository interface {
	Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	GetByID(ctx context.Context, id uint) (*domain.Movie, error)
	GetAll(ctx context.Context) ([]domain.Movie, error)
	Update(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id uint) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *movieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(movie).Error; err != nil {
			return ErrCreateMovie
		}

		// If the movie has related entities (e.g., genres, actors),we can save them here
		// Example: tx.Create(&movie.Actors)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *movieRepository) GetByID(ctx context.Context, id uint) (*domain.Movie, error) {
	var movie domain.Movie
	if err := r.db.WithContext(ctx).First(&movie, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	return &movie, nil
}

func (r *movieRepository) GetAll(ctx context.Context) ([]domain.Movie, error) {
	var movies []domain.Movie
	if err := r.db.WithContext(ctx).Find(&movies).Error; err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	return movies, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(movie).Error; err != nil {
			return ErrUpdateMovie
		}

		// If we need to update related entities, we can do it here
		// Example: tx.Model(&movie).Association("Actors").Replace(movie.Actors)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *movieRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Movie{}, id).Error
}
