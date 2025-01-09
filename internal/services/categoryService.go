// category service

package services

import (
	"errors"

	models "github.com/mackenzii/freemusic/internal/models"
	"gorm.io/gorm"
)

// CategoryService provides services for managing categories
type CategoryService struct {
	DB *gorm.DB
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

// CreateCategory creates a new category in the database
func (s *CategoryService) CreateCategory(category *models.Category) error {
	if err := s.DB.Create(category).Error; err != nil {
		return err
	}

	return nil
}

// GetCategoryByID retrieves a category by its ID
func (s *CategoryService) GetCategoryByID(categoryID int) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory updates an existing category in the database
func (s *CategoryService) UpdateCategory(categoryID int, category *models.Category) (*models.Category, error) {
	if err := s.DB.Where("category_id = ?", categoryID).Updates(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category from the database
func (s *CategoryService) DeleteCategory(categoryID int) error {
	if err := s.DB.Where("category_id = ?", categoryID).Delete(&models.Category{}).Error; err != nil {
		return err
	}

	return nil
}
