package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mackenzii/freemusic/internal/config"
	"github.com/mackenzii/freemusic/internal/database"
	"github.com/mackenzii/freemusic/internal/fixtures"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/mackenzii/freemusic/internal/server"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env

	err := godotenv.Load("/root/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Charger la configuration depuis l'environnement
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialisation de la base de données avec la configuration
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	//migration
	// Effectuer la migration
	err = db.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatal("failed to migrate schema", err)
	}

	// Appliquer les migrations pour créer les tables dans la base de données
	err = db.AutoMigrate(

		&models.Users{},
		&models.Artist{},
		&models.Category{},
		&models.Event{},
	)
	if err != nil {
		log.Fatal("Erreur lors de la migration des modèles :", err)
	}

	// Génération des fixtures
	log.Println("Generating fixtures...")
	fixtures.GenerateUsers(db)
	fixtures.GenerateArtists(db)
	fixtures.GenerateCategories(db)
	fixtures.GenerateEvents(db)
	log.Println("Fixtures generated successfully!")

	// Démarrage du serveur
	server.Run()
}
