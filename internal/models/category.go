package models

type Category struct {
	CategoryID int    `gorm:"primaryKey;autoIncrement" json:"category_id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
}
