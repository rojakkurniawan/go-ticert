package repository

import (
	"ticert/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *entity.Category) error
	GetCategoryByID(id uuid.UUID) (*entity.Category, error)
	GetCategories(eventID uuid.UUID) ([]*entity.Category, error)
	UpdateCategory(category *entity.Category) error
	DeleteCategory(id uuid.UUID) error
	CheckStock(categoryID uuid.UUID, quantity int) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *entity.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.Where("id = ?", id).Preload("Event").First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategories(eventID uuid.UUID) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := r.db.Where("event_id = ?", eventID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) UpdateCategory(category *entity.Category) error {
	if err := r.db.Model(&entity.Category{}).Where("id = ?", category.ID).Updates(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(id uuid.UUID) error {
	if err := r.db.Delete(&entity.Category{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) CheckStock(categoryID uuid.UUID, quantity int) (bool, error) {
	var category entity.Category
	if err := r.db.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return false, err
	}
	return category.Quantity >= quantity, nil
}
