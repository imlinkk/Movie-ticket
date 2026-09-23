package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketStatus string

const (
	TicketStatusBooked    TicketStatus = "booked"
	TicketStatusPaid      TicketStatus = "paid"
	TicketStatusCancelled TicketStatus = "cancelled"
)

type Ticket struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShowtimeID  uint           `gorm:"not null;index;uniqueIndex:idx_showtime_seat,priority:1" json:"showtime_id"`
	Showtime    Showtime       `gorm:"foreignKey:ShowtimeID" json:"showtime,omitempty"`
	SeatNumber  string         `gorm:"size:10;not null;uniqueIndex:idx_showtime_seat,priority:2" json:"seat_number"`
	Price       float64        `gorm:"not null" json:"price"`
	Status      TicketStatus   `gorm:"size:20;default:'booked';not null" json:"status"`
	BookingTime time.Time      `json:"booking_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type BookTicketRequest struct {
	ShowtimeID uint   `json:"showtime_id" binding:"required"`
	SeatNumber string `json:"seat_number" binding:"required,min=1,max=10"`
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=booked paid cancelled"`
}
