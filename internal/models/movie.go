package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	DurationMinutes int            `gorm:"not null" json:"duration_minutes"`
	Genre           string         `gorm:"size:100;not null;index" json:"genre"`
	ReleaseDate     time.Time      `json:"release_date"`
	PosterURL       string         `gorm:"size:500" json:"poster_url"`
	TrailerURL      string         `gorm:"size:500" json:"trailer_url"`
	Rating          float64        `gorm:"default:0.0" json:"rating"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateMovieRequest struct {
	Title           string    `json:"title" binding:"required,min=1,max=200"`
	Description     string    `json:"description" binding:"omitempty"`
	DurationMinutes int       `json:"duration_minutes" binding:"required,gt=0"`
	Genre           string    `json:"genre" binding:"required,min=1,max=100"`
	ReleaseDate     time.Time `json:"release_date" binding:"required"`
	PosterURL       string    `json:"poster_url" binding:"omitempty,url"`
	TrailerURL      string    `json:"trailer_url" binding:"omitempty,url"`
	Rating          float64   `json:"rating" binding:"omitempty,gte=0,lte=10"`
}

type UpdateMovieRequest struct {
	Title           string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description     string     `json:"description" binding:"omitempty"`
	DurationMinutes int        `json:"duration_minutes" binding:"omitempty,gt=0"`
	Genre           string     `json:"genre" binding:"omitempty,min=1,max=100"`
	ReleaseDate     *time.Time `json:"release_date" binding:"omitempty"`
	PosterURL       string     `json:"poster_url" binding:"omitempty,url"`
	TrailerURL      string     `json:"trailer_url" binding:"omitempty,url"`
	Rating          *float64   `json:"rating" binding:"omitempty,gte=0,lte=10"`
}
