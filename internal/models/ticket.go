package models

import "time"

type Ticket struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint `gorm:"not null"`
	User          Users
	EventID       uint `gorm:"not null"`
	Event         Event
	PurchaseDate  time.Time `gorm:"not null"`
	SeatNumber    *string
	ReservationID uint
}
