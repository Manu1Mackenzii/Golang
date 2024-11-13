package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mackenzii/freemusic/internal/database"
	"github.com/mackenzii/freemusic/internal/models"
)

type TicketController struct{}

// GetTickets retrieves all tickets
func (tc *TicketController) GetTickets(c *fiber.Ctx) error {
	var tickets []models.Ticket
	if err := database.DB.Find(&tickets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tickets",
		})
	}
	return c.JSON(tickets)
}

// CreateTicket creates a new ticket
func (tc *TicketController) CreateTicket(c *fiber.Ctx) error {
	var ticket models.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if err := database.DB.Create(&ticket).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create ticket",
		})
	}
	return c.JSON(ticket)
}

// DeleteTicket deletes a ticket by ID
func (tc *TicketController) DeleteTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.Ticket{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete ticket",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
