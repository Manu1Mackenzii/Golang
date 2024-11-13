package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mackenzii/freemusic/internal/controllers"
	"github.com/mackenzii/freemusic/internal/services"
	"gorm.io/gorm"
)

// EventRoutes defines the routes for events
func EventRoutes(router *gin.Engine, db *gorm.DB) {
	eventService := services.NewEventService(db)
	eventController := controllers.NewEventController(eventService)

	router.GET("/events", eventController.GetAllEvents)
	router.GET("/events/:event_id", eventController.GetEventByID)
	router.POST("/events", eventController.CreateEvent)
	router.PUT("/events/:event_id", eventController.UpdateEvent)
	router.DELETE("/events/:event_id", eventController.DeleteEvent)
}
