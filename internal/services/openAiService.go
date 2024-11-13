package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mackenzii/freemusic/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService() *OpenAIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client: client,
	}
}

// SuggestConcerts suggère des concerts ou des événements en fonction des préférences musicales des utilisateurs
func (s *OpenAIService) SuggestConcerts(users []models.Users) ([]string, error) {
	var suggestions []string

	// Préparation du prompt pour OpenAI basé sur les préférences musicales des utilisateurs
	prompt := "Voici les préférences musicales de certains utilisateurs :\n"
	for _, user := range users {
		prompt += fmt.Sprintf("Utilisateur: %s, Genre préféré: %s, Artiste préféré: %s, Niveau de compétence: %s, Biographie: %s\n",
			user.Username, user.FavoriteSport, user.FavoriteSport, user.SkillLevel, user.Bio)
	}
	prompt += "Sur la base de ces informations, suggère des concerts ou événements musicaux adaptés, ainsi que des artistes à suivre. Par exemple : 'Concert de l'artiste XYZ à Paris le 20 janvier'."

	// Création de la requête à l'API OpenAI pour obtenir des suggestions d'événements
	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-4", // Utilisation du modèle GPT-4 pour des résultats plus précis
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a music event and concert expert.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		log.Printf("Erreur lors de l'appel à l'API OpenAI: %v", err)
		return nil, err
	}

	// Extraction des suggestions d'événements ou de concerts depuis la réponse d'OpenAI
	if len(resp.Choices) > 0 {
		suggestions = append(suggestions, resp.Choices[0].Message.Content)
	} else {
		return nil, fmt.Errorf("aucune suggestion d'événements générée par OpenAI")
	}

	return suggestions, nil
}
