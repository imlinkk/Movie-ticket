package models

import (
	"time"

	"gorm.io/gorm"
)

type Showtime struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MovieID        uint           `gorm:"not null;index" json:"movie_id"`
	Movie          Movie          `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	CinemaID       uint           `gorm:"not null;index" json:"cinema_id"`
	Cinema         Cinema         `gorm:"foreignKey:CinemaID" json:"cinema,omitempty"`
	RoomName       string         `gorm:"size:50;not null" json:"room_name"`
	StartTime      time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime        time.Time      `gorm:"not null" json:"end_time"`
	Price          float64        `gorm:"not null" json:"price"`
	TotalSeats     int            `gorm:"not null" json:"total_seats"`
	AvailableSeats int            `gorm:"not null" json:"available_seats"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateShowtimeRequest struct {
	MovieID    uint      `json:"movie_id" binding:"required"`
	CinemaID   uint      `json:"cinema_id" binding:"required"`
	RoomName   string    `json:"room_name" binding:"required,min=1,max=50"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	Price      float64   `json:"price" binding:"required,gt=0"`
	TotalSeats int       `json:"total_seats" binding:"required,gt=0"`
}

type UpdateShowtimeRequest struct {
	RoomName  string     `json:"room_name" binding:"omitempty,min=1,max=50"`
	StartTime *time.Time `json:"start_time" binding:"omitempty"`
	Price     *float64   `json:"price" binding:"omitempty,gt=0"`
}
