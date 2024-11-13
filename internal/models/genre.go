package models

type Genre struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"size:100;not null;unique"`
	Events []Event `gorm:"many2many:event_genres;"`
}
