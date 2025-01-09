package models

import (
	"encoding/json"
)

// Artist représente la table artist dans la base de données.
type Artist struct {
	ArtistID    int             `json:"artist_id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"type:varchar(255)"`
	Bio         string          `json:"bio" gorm:"type:text"`
	SocialLinks json.RawMessage `json:"social_links" gorm:"type:jsonb"`
}
