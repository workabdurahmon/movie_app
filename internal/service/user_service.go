package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"movie_app/internal/domain"
	"movie_app/internal/repository"
)

const (
	TokenExpirationHours = 24
)

var (
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrPasswordHashFailed = errors.New("failed to hash password")
)

type UserService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	GetUser(ctx context.Context, id uint) (*domain.User, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) *userService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	existing, _ := s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("checking password: %w", ErrInvalidCredentials)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(TokenExpirationHours * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("signing token: %w", err)
	}

	return &domain.LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user by ID: %w", err)
	}

	return result, nil
}
