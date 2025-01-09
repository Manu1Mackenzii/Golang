package fixtures

import (
	"log"

	"github.com/mackenzii/freemusic/internal/models"
	"gorm.io/gorm"
)

func GenerateCategories(db *gorm.DB) error {
	// Liste de titres de concert et icônes associés
	titles := []string{
		"Concert Rock Extravaganza",
		"Festival de Jazz Étoilé",
		"Musique Classique en Plein Air",
		"Rythmes Électro de Nuit",
		"Concert Pop Moderne",
		"Festival Acoustique Évasion",
		"Sounds of the 80s",
		"Concert de Musique Latine",
		"Indie Vibes Showcase",
		"Concert de Musique du Monde",
	}

	icons := []string{
		"(🎸)", "(🎷)", "(🎻)", "(🎧)", "(🎤)", "(🎶)", "(📀)", "(🎹)", "(🎺)", "(🌍)",
	}

	// Exemple de génération de catégories
	for i := 0; i < len(titles); i++ {
		category := models.Category{
			CategoryID: i + 1,
			Name:       titles[i],
			Icon:       icons[i],
		}

		// Vérification si la catégorie existe déjà avant insertion
		var existingCategory models.Category
		if err := db.First(&existingCategory, "category_id = ?", category.CategoryID).Error; err == nil {
			log.Printf("La catégorie %d existe déjà, elle sera ignorée.\n", category.CategoryID)
			continue
		}

		// Insertion de la catégorie si elle n'existe pas
		if err := db.Create(&category).Error; err != nil {
			log.Printf("Erreur lors de l'insertion de la catégorie %d: %v\n", category.CategoryID, err)
			return err
		}
	}

	return nil
}
