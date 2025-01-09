package fixtures

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mackenzii/freemusic/internal/models"

	"gorm.io/gorm"
)

// GenerateArtists insère des données pour les artistes
func GenerateArtists(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		artist := models.Artist{
			Name: fmt.Sprintf("Artist %d", i+1),
			Bio:  fmt.Sprintf("Bio for artist %d", i+1),
			SocialLinks: json.RawMessage(fmt.Sprintf(`{
				"facebook": "https://facebook.com/artist%d",
				"twitter": "https://twitter.com/artist%d"
			}`, i+1, i+1)),
		}
		if err := db.Create(&artist).Error; err != nil {
			log.Printf("Failed to create artist %d: %v", i+1, err)
		}
	}
	log.Println("Artists fixtures loaded successfully!")
}
