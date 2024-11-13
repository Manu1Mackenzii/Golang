package controllers

import (
	"github.com/mackenzii/freemusic/internal/database"

	"github.com/mackenzii/freemusic/internal/models"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

// GetUser retrieves user information by ID
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Users

	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
}

// CreateUser creates a new user
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.JSON(user)
}

// DeleteUser deletes a user by ID
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.Users{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
