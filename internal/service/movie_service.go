package service

import (
	"context"
	"errors"
	"fmt"
	"movie_app/internal/domain"
	"movie_app/internal/repository"
	"time"
)

var (
	ErrInvalidYear     = errors.New("invalid year")
	ErrInvalidDuration = errors.New("invalid duration")
	ErrInvalidRating   = errors.New("rating must be between 1 and 10")
	ErrMovieNotFound   = errors.New("movie not found")
)

type MovieService interface {
	Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	GetByID(ctx context.Context, id uint) (*domain.Movie, error)
	GetAll(ctx context.Context) ([]domain.Movie, error)
	Update(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id uint) error
}

type movieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *movieService {
	return &movieService{
		repo: repo,
	}
}

func (s *movieService) Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if err := s.ValidateMovie(movie); err != nil {
		return nil, fmt.Errorf("validating movie: %w", err)
	}

	result, err := s.repo.Create(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("creating movie: %w", err)
	}

	return result, nil
}

func (s *movieService) GetByID(ctx context.Context, id uint) (*domain.Movie, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting movie by ID: %w", err)
	}

	return result, nil
}

func (s *movieService) GetAll(ctx context.Context) ([]domain.Movie, error) {
	result, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all movies: %w", err)
	}

	return result, nil
}

func (s *movieService) Update(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if err := s.ValidateMovie(movie); err != nil {
		return nil, fmt.Errorf("validating movie: %w", err)
	}

	result, err := s.repo.Update(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("updating movie: %w", err)
	}

	return result, nil
}

func (s *movieService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting movie: %w", err)
	}

	return nil
}

func (s *movieService) ValidateMovie(movie *domain.Movie) error {
	currentYear := time.Now().Year()
	if movie.Year < 1888 || movie.Year > currentYear+5 {
		return ErrInvalidYear
	}

	if movie.Duration <= 0 {
		return ErrInvalidDuration
	}

	if movie.Rating < 1 || movie.Rating > 10 {
		return ErrInvalidRating
	}

	return nil
}
