package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	JWT struct {
		Secret          string
		ExpirationHours int
	}
}

func LoadConfig() (*Config, error) {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load("/root/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configurer viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app") // Chemin où le fichier de config.yaml pourrait être copié

	// Charger la configuration YAML
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Remplacer les variables d'environnement dans le fichier YAML
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
