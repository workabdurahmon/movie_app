package domain

import (
	"time"
)

type Movie struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Director    string    `json:"director" gorm:"not null"`
	Year        int       `json:"year" gorm:"not null"`
	Plot        string    `json:"plot" gorm:"type:text"`
	Genre       string    `json:"genre" gorm:"not null"`
	Rating      float64   `json:"rating" gorm:"type:decimal(2,1)"`
	Duration    int       `json:"duration" gorm:"not null"` // Duration in minutes
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateMovieRequest struct {
	Title       string    `json:"title" binding:"required"`
	Director    string    `json:"director" binding:"required"`
	Year        int       `json:"year" binding:"required"`
	Plot        string    `json:"plot"`
	Genre       string    `json:"genre" binding:"required"`
	Rating      float64   `json:"rating" binding:"required"`
	Duration    int       `json:"duration" binding:"required"`
}

type UpdateMovieRequest struct {
	Title       *string    `json:"title"`
	Director    *string    `json:"director"`
	Year        *int       `json:"year"`
	Plot        *string    `json:"plot"`
	Genre       *string    `json:"genre"`
	Rating      *float64   `json:"rating"`
	Duration    *int       `json:"duration"`
}
