package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	Upcoming  Status = "upcoming"
	Ongoing   Status = "ongoing"
	Completed Status = "completed"
	Expired   Status = "expired"
)

type Event struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64          `gorm:"not null"`
	Title         string         `gorm:"null"`
	Description   string         `gorm:"null"`
	LocationID    string         `gorm:"null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	EventDate     time.Time      `gorm:"null"`
	EventTime     time.Time      `gorm:"null"`
	EndTime       time.Time      `gorm:"null"`
	Address       string         `gorm:"null"`
	Latitude      float64        `gorm:"null"`
	Longitude     float64        `gorm:"null"`
	Status        Status         `gorm:"null"`
	ArtistID      int            `gorm:"null"`
	CategoryIds   string         `gorm:"type:json"` // Stockage des catégories sous forme de JSON
	GalleryImages string         `gorm:"type:json"` // Stockage des images sous forme de JSON
}
