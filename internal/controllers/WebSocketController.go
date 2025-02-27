package controllers

import (
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/mackenzii/freemusic/internal/services"
)

type WebSocketController struct {
	service *services.WebSocketService
}

func NewWebSocketController(service *services.WebSocketService) *WebSocketController {
	return &WebSocketController{service: service}
}

func (ctrl *WebSocketController) WebSocketHandler(c *websocket.Conn) {
	log.Println("New WebSocket connection established")
	defer func() {
		log.Println("WebSocket connection closed")
	}()

	ctrl.service.HandleWebSocket(c)
}
