package models

type Genre struct {
	ID     int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string  `gorm:"size:100;not null;unique"`
	Events []Event `gorm:"many2many:event_genres;"`
}
