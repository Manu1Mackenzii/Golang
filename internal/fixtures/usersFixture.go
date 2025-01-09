package fixtures

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mackenzii/freemusic/internal/models"
	"gorm.io/gorm"
)

func GenerateUsers(db *gorm.DB) error {
	// Liste des rôles possibles
	roles := []models.Role{
		models.RoleAdmin,
		models.RoleOrganizer,
		models.RoleUser,
	}

	// Générer 100 utilisateurs avec des rôles aléatoires
	for i := 1; i <= 10; i++ {
		// Choisir un rôle aléatoire
		role := roles[rand.Intn(len(roles))]
		if i%10 == 0 { // Par exemple, attribuer le rôle "organizer" à tous les 10e utilisateur
			role = models.RoleOrganizer
		}

		user := models.Users{
			ID:           fmt.Sprintf("user_%03d", i),
			Username:     fmt.Sprintf("username_%03d", i),
			Email:        fmt.Sprintf("user%03d@example.com", i),
			PasswordHash: "hashed_password", // Remplacer par un mot de passe haché approprié
			UpdatedAt:    time.Now(),
			Role:         role,
			Location:     "Paris",
			Latitude:     48.8566 + float64(i)*0.0001,
			Longitude:    2.3522 + float64(i)*0.0001,
			SkillLevel:   "beginner", // Exemple de niveau
			Bio:          "Ceci est un exemple de bio.",
			FCMToken:     fmt.Sprintf("fcm_token_%03d", i),
		}

		// Vérifier si l'utilisateur existe déjà ou le créer
		if err := db.FirstOrCreate(&user, "id = ?", user.ID).Error; err != nil {
			log.Printf("Erreur lors de l'insertion ou de la recherche de l'utilisateur %s: %v\n", user.ID, err)
			return err
		}
		log.Printf("Utilisateur %s avec rôle %s ajouté ou existant.\n", user.ID, user.Role)
	}

	return nil
}
