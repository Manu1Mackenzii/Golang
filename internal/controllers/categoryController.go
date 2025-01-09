package controllers

import (
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/mackenzii/freemusic/internal/services"
	"gorm.io/gorm"
)

type CategoryController struct {
	CategoryService *services.CategoryService
	RedisClient     *redis.Client
	AuthService     *services.AuthService
	DB              *gorm.DB
}

// NewCategoryController crée une nouvelle instance de CategoryController
func NewCategoryController(categoryService *services.CategoryService, authService *services.AuthService, db *gorm.DB, redisClient *redis.Client) *CategoryController {
	return &CategoryController{
		CategoryService: categoryService,
		AuthService:     authService,
		DB:              db,
		RedisClient:     redisClient,
	}
}

// GetCategories retourne toutes les catégories
func (cc *CategoryController) GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := cc.DB.Find(&categories).Error; err != nil {
		log.Printf("Échec de la récupération des catégories : %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Échec de la récupération des catégories"})
	}

	return c.JSON(categories)
}

// GetCategory retourne une catégorie spécifique
func (c *CategoryController) GetCategory(ctx *fiber.Ctx) error {
	categoryIdStr := ctx.Params("id") // Récupérer l'ID depuis les paramètres d'URL

	// Convertir categoryIdStr en int
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		// Si la conversion échoue, retourner une erreur 400 (mauvaise requête)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID de catégorie invalide",
		})
	}

	category, err := c.CategoryService.GetCategoryByID(categoryId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Catégorie non trouvée",
		})
	}

	return ctx.JSON(category)
}

// CreateCategory crée une nouvelle catégorie
func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		log.Printf("Échec du binding de la catégorie : %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Requête invalide"})
	}

	if err := cc.DB.Create(&category).Error; err != nil {
		log.Printf("Échec de la création de la catégorie : %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Échec de la création de la catégorie"})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategory met à jour une catégorie existante
func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Category
	if err := cc.DB.First(&category, "category_id = ?", id).Error; err != nil {
		log.Printf("Échec de la récupération de la catégorie %s : %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Catégorie non trouvée"})
	}

	if err := c.BodyParser(&category); err != nil {
		log.Printf("Échec du binding de la catégorie : %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Requête invalide"})
	}

	if err := cc.DB.Save(&category).Error; err != nil {
		log.Printf("Échec de la mise à jour de la catégorie %s : %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Échec de la mise à jour de la catégorie"})
	}

	return c.JSON(category)
}

// DeleteCategory supprime une catégorie
func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Category
	if err := cc.DB.First(&category, "category_id = ?", id).Error; err != nil {
		log.Printf("Échec de la récupération de la catégorie %s : %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Catégorie non trouvée"})
	}

	if err := cc.DB.Delete(&category).Error; err != nil {
		log.Printf("Échec de la suppression de la catégorie %s : %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Échec de la suppression de la catégorie"})
	}

	return c.JSON(fiber.Map{"message": "Catégorie supprimée avec succès"})
}
