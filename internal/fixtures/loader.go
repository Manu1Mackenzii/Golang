package fixtures

import (
	"log"

	"gorm.io/gorm"
)

// LoadFixtures charge toutes les fixtures
func LoadFixtures(db *gorm.DB) {
	log.Println("Loading fixtures...")
	GenerateUsers(db)      // Charger les utilisateurs
	GenerateArtists(db)    // Charger les artistes
	GenerateCategories(db) // Charger les catégories
	GenerateEvents(db)     // Charger les événements

	log.Println("All fixtures loaded successfully!")
}
