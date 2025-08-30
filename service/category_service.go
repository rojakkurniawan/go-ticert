package service

import (
	"context"
	"errors"
	"ticert/dto/request"
	"ticert/dto/response"
	"ticert/entity"
	"ticert/repository"
	"ticert/utils/errs"
	"ticert/utils/validator"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*response.CategoryResponse, map[string]string, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*response.CategoryResponse, error)
	GetCategories(ctx context.Context, eventID uuid.UUID) ([]*response.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *request.UpdateCategoryRequest) (*response.CategoryResponse, map[string]string, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	eventRepo    repository.EventRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, eventRepo repository.EventRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo, eventRepo: eventRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*response.CategoryResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	_, err := s.eventRepo.GetEventByID(req.EventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errs.ErrEventNotFound
		}
		return nil, nil, errs.ErrInternalServerError
	}

	parsedDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, map[string]string{"event_date": "Invalid event date format. Use YYYY-MM-DD"}, nil
	}

	category := &entity.Category{
		EventID:   req.EventID,
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		EventDate: parsedDate,
	}

	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewCategoryResponse(category), nil, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, errs.ErrInternalServerError
	}

	return response.NewCategoryResponse(category), nil
}

func (s *categoryService) GetCategories(ctx context.Context, eventID uuid.UUID) ([]*response.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetCategories(eventID)
	if err != nil {
		return nil, errs.ErrInternalServerError
	}

	if len(categories) == 0 {
		return nil, errs.ErrCategoryNotFound
	}

	var categoryResponses []*response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.NewCategoryResponse(category))
	}

	return categoryResponses, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req *request.UpdateCategoryRequest) (*response.CategoryResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	parsedDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, map[string]string{"event_date": "Invalid event date format. Use YYYY-MM-DD"}, nil
	}

	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errs.ErrCategoryNotFound
		}
		return nil, nil, errs.ErrInternalServerError
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Price != 0 {
		category.Price = req.Price
	}

	if req.Quantity != 0 {
		category.Quantity = req.Quantity
	}

	if req.EventDate != "" {
		category.EventDate = parsedDate
	}

	if err := s.categoryRepo.UpdateCategory(category); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewCategoryResponse(category), nil, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrCategoryNotFound
		}
		return errs.ErrInternalServerError
	}
	if err := s.categoryRepo.DeleteCategory(id); err != nil {
		return errs.ErrInternalServerError
	}
	return nil
}
