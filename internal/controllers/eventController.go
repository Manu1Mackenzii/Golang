package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/mackenzii/freemusic/internal/services"
	"gorm.io/gorm"
)

type EventController struct {
	EventService *services.EventService
	ChatService  *services.ChatService
	RedisClient  *redis.Client
	AuthService  *services.AuthService
	DB           *gorm.DB
}

type GeoResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

// NewEventController creates a new EventController instance
func NewEventController(eventService *services.EventService, authService *services.AuthService, db *gorm.DB, chatService *services.ChatService, redisClient *redis.Client) *EventController {
	return &EventController{
		EventService: eventService,
		AuthService:  authService,
		DB:           db,
		ChatService:  chatService,
		RedisClient:  redisClient,
	}
}

// CreateEvent handles the creation of a new event
func (ec *EventController) CreateEvent(c *fiber.Ctx) error {
	var event models.Event

	// Parse input data
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate user authentication
	userID, err := ec.AuthService.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	event.UserID = userID

	// Validate and fetch coordinates for the address
	lat, lng, err := getCoordinates(event.Address, "YOUR_API_KEY")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid address"})
	}
	event.Latitude = lat
	event.Longitude = lng

	// Create the event
	createdEvent, err := ec.EventService.CreateEvent(&event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create event"})
	}

	return c.JSON(createdEvent)
}

// GetAllEvents retrieves all events
func (ec *EventController) GetAllEvents(c *fiber.Ctx) error {
	events, err := ec.EventService.GetAllEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch events"})
	}
	return c.JSON(events)
}

// GetEventByID retrieves an event by its ID
func (ec *EventController) GetEventByID(c *fiber.Ctx) error {
	eventID := c.Params("event_id")

	// Convert event ID to int
	eventIDInt, err := strconv.Atoi(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	event, err := ec.EventService.GetEventByID(eventIDInt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}
	return c.JSON(event)
}

// UpdateEvent updates an existing event
func (ec *EventController) UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("event_id")
	var req models.Event

	// Parse input data
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	updatedEvent, err := ec.EventService.UpdateEvent(eventID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating event"})
	}
	return c.JSON(updatedEvent)
}

// DeleteEvent deletes an event
func (ec *EventController) DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("event_id")

	// Convert event ID to int
	eventIDInt, err := strconv.Atoi(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	if err := ec.EventService.DeleteEvent(eventIDInt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting event"})
	}
	return c.JSON(fiber.Map{"message": "Event deleted successfully"})
}

// getCoordinates fetches coordinates for a given address using Google Geocoding API
func getCoordinates(address, apiKey string) (float64, float64, error) {
	address = url.QueryEscape(address)
	apiURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var geoResp GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, err
	}

	if len(geoResp.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	lat := geoResp.Results[0].Geometry.Location.Lat
	lng := geoResp.Results[0].Geometry.Location.Lng
	return lat, lng, nil
}
