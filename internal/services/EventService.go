package services

import (
	"errors"
	"time"

	models "github.com/mackenzii/freemusic/internal/models"
	"gorm.io/gorm"
)

// EventService provides services for managing events
type EventService struct {
	DB *gorm.DB
}

// NewEventService creates a new instance of EventService
func NewEventService(db *gorm.DB) *EventService {
	return &EventService{DB: db}
}

// CreateEvent creates a new event in the database
func (s *EventService) CreateEvent(event *models.Event, UserID int64) error {
	event.UserID = UserID
	event.Status = "upcoming"
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := s.DB.Create(event).Error; err != nil {
		return err
	}

	return nil
}

// GetEventByID retrieves an event by its ID
func (s *EventService) GetEventByID(eventID int) (*models.Event, error) {
	var event models.Event
	if err := s.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

// GetAllEvents retrieves all events
func (s *EventService) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	if err := s.DB.Where("deleted_at IS NULL").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates an existing event in the database
func (s *EventService) UpdateEvent(eventID string, event *models.Event) (*models.Event, error) {
	event.UpdatedAt = time.Now()

	if err := s.DB.Save(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

// DeleteEvent soft deletes an event
func (s *EventService) DeleteEvent(eventID int) error {
	var event models.Event
	if err := s.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}

	if err := s.DB.Save(&event).Error; err != nil {
		return err
	}

	return nil
}

// UpdateEventStatus updates the status of an event (e.g., "ongoing", "completed")
func (s *EventService) UpdateEventStatus(eventID int, status string) error {
	var event models.Event
	if err := s.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}

	// Validate status
	if status != "upcoming" && status != "ongoing" && status != "completed" {
		return errors.New("invalid event status")
	}

	event.UpdatedAt = time.Now()

	if err := s.DB.Save(&event).Error; err != nil {
		return err
	}

	return nil
}

// GetEventsByOrganizerID retrieves events by the organizer's ID
func (s *EventService) GetEventsByOrganizerID(organizerID int) ([]models.Event, error) {
	var events []models.Event
	if err := s.DB.Where("organizer_id = ? AND deleted_at IS NULL", organizerID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// SoftDeleteEvent allows an event to be marked as deleted without actually removing it from the database
func (s *EventService) SoftDeleteEvent(eventID int) error {
	var event models.Event
	if err := s.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}

	if err := s.DB.Save(&event).Error; err != nil {
		return err
	}

	return nil
}
