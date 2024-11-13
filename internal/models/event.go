package models

import "time"

type Status string

const (
	Upcoming  Status = "upcoming"
	Ongoing   Status = "ongoing"
	Completed Status = "completed"
	Expired   Status = "expired"
)

// Event represents the event table structure
type Event struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	UserID      int        `json:"user_id"`
	OrganizerID string     `json:"organizer_id" gorm:"not null"` // Référence vers l'ID de l'organisateur (Users.id)
	Organizer   Users      `json:"organizer" gorm:"foreignKey:OrganizerID"`
	EventID     string     `gorm:"primaryKey;column:event_id"` // Changé en string pour utiliser un ulid
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	LocationID  string     `gorm:"column:location_id"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;default:null"`
	EventDate   time.Time  `gorm:"column:event_date"`
	EventTime   time.Time  `gorm:"column:event_time"`
	EndTime     time.Time  `gorm:"column:end_time"`
	Address     string     `gorm:"column:address"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Status      Status     `json:"status"` // Statut de l'événement

}
