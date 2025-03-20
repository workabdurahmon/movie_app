package repository

import (
	"context"
	"errors"
	"movie_app/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrCreateUser = errors.New("failed to create user")
	ErrUpdateUser = errors.New("failed to update user")
	ErrFetchUser  = errors.New("failed to fetch user")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, ErrCreateUser
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
