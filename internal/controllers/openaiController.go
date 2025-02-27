package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mackenzii/freemusic/internal/services"
)

type OpenAiController struct {
	openAi              *services.OpenAIService
	matchPlayersService *services.MatchPlayersService
}

func NewOpenAiController(openAi *services.OpenAIService, matchPlayersService *services.MatchPlayersService) *OpenAiController {
	return &OpenAiController{
		openAi:              openAi,
		matchPlayersService: matchPlayersService,
	}
}

// @Summary GetFormationFromAi
// @Description Get a suggested formation from the AI
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param match_id path string true "Match ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/openai/formation/{match_id} [get]
func (ctrl *OpenAiController) GetFormationFromAi(c *fiber.Ctx) error {
	// matchID := c.Params("match_id")

	// Récupérer les joueurs du match
	// matchPlayers, err := ctrl.matchPlayersService.GetMatchPlayersByMatchID(matchID)
	// if err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Match players not found for match ID: " + matchID})
	// }

	// var players []models.Users
	// for _, matchPlayer := range matchPlayers {
	// 	players = append(players, matchPlayer.Player)
	// }

	// Suggérer une formation basée sur les statistiques des joueurs
	// formation, err := ctrl.openAi.SuggestFormations(players)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
	// 		"error": err.Error(),
	// 	})
	// }
	return c.JSON(map[string]interface{}{
		// 	"formation": formation,
	})
}
