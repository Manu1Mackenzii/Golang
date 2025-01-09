package fixtures

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/mackenzii/freemusic/internal/models"
	"gorm.io/gorm"
)

// Liste des statuts possibles
var eventStatuses = []models.Status{
	models.Upcoming,  // Événement à venir
	models.Ongoing,   // Événement en cours
	models.Completed, // Événement terminé
}

func GenerateEvents(db *gorm.DB) error {
	for i := 1; i <= 10; i++ {
		// Sélectionner un utilisateur existant
		var user models.Users
		if err := db.First(&user).Error; err != nil {
			return fmt.Errorf("utilisateur non trouvé : %v", err)
		}
		// Convertir l'ID utilisateur (string) en int64
		userID, err := strconv.ParseInt(user.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("erreur lors de la conversion de l'ID utilisateur : %v", err)
		}

		// Créer l'événement
		event := models.Event{
			ID:            int64(i),
			UserID:        userID, // Assignation de l'ID utilisateur existant
			Title:         fmt.Sprintf("Event %d", i),
			Description:   fmt.Sprintf("Description of event %d", i),
			LocationID:    fmt.Sprintf("loc%03d", i),
			EventDate:     time.Now().Add(time.Duration(i) * time.Hour),
			EventTime:     time.Now().Add(time.Duration(i) * time.Hour),
			EndTime:       time.Now().Add(time.Duration(i+1) * time.Hour),
			Address:       fmt.Sprintf("Address %d", i),
			Latitude:      48.8566, // Exemple de latitude
			Longitude:     2.3522,  // Exemple de longitude
			ArtistID:      i,
			CategoryIds:   fmt.Sprintf("[\"Category %d\"]", i),
			GalleryImages: fmt.Sprintf("[\"image%d.jpg\"]", i),
		}

		// Attribuer un statut aléatoire
		event.Status = eventStatuses[rand.Intn(len(eventStatuses))] // Sélection aléatoire dans la liste des statuts

		// Vérification si l'événement existe déjà
		var existingEvent models.Event
		if err := db.First(&existingEvent, "id = ?", event.ID).Error; err == nil {
			log.Printf("L'événement %d existe déjà, il sera ignoré.\n", event.ID)
			continue
		}

		// Insertion de l'événement si il n'existe pas
		if err := db.Create(&event).Error; err != nil {
			log.Printf("Erreur lors de l'insertion de l'événement %d: %v\n", event.ID, err)
			return err
		}
	}

	return nil
}
