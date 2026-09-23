package models

import (
	"time"

	"gorm.io/gorm"
)

type Cinema struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:150;not null" json:"name"`
	Address    string         `gorm:"size:255;not null" json:"address"`
	City       string         `gorm:"size:100;not null;index" json:"city"`
	TotalRooms int            `gorm:"default:1" json:"total_rooms"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateCinemaRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=150"`
	Address    string `json:"address" binding:"required,min=5,max=255"`
	City       string `json:"city" binding:"required,min=2,max=100"`
	TotalRooms int    `json:"total_rooms" binding:"required,gt=0"`
}

type UpdateCinemaRequest struct {
	Name       string `json:"name" binding:"omitempty,min=2,max=150"`
	Address    string `json:"address" binding:"omitempty,min=5,max=255"`
	City       string `json:"city" binding:"omitempty,min=2,max=100"`
	TotalRooms int    `json:"total_rooms" binding:"omitempty,gt=0"`
}
