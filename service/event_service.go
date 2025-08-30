package service

import (
	"context"
	"errors"
	"ticert/dto/request"
	"ticert/dto/response"
	"ticert/entity"
	"ticert/repository"
	"ticert/utils/errs"
	utils_response "ticert/utils/response"
	"ticert/utils/validator"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(ctx context.Context, req *request.CreateEventRequest) (*response.EventResponse, map[string]string, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*response.EventResponse, error)
	GetEvents(ctx context.Context, req *request.GetEventsRequest) (*response.EventListResponse, map[string]string, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, req *request.UpdateEventRequest) (*response.EventResponse, map[string]string, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) CreateEvent(ctx context.Context, req *request.CreateEventRequest) (*response.EventResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	title, _ := s.eventRepo.GetEventByTitle(req.Title)
	if title != nil {
		return nil, nil, errs.ErrEventTitleAlreadyExists
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, map[string]string{"start_date": "Invalid start date format. Use YYYY-MM-DD"}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, map[string]string{"end_date": "Invalid end date format. Use YYYY-MM-DD"}, nil
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, map[string]string{"start_time": "Invalid start time format. Use HH:MM"}, nil
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, map[string]string{"end_time": "Invalid end time format. Use HH:MM"}, nil
	}

	if startDate.After(endDate) {
		return nil, map[string]string{"start_date": "Start date must be before end date"}, nil
	}

	if endDate.Before(startDate) {
		return nil, map[string]string{"end_date": "End date must be after start date"}, nil
	}

	baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	startDateTime := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
	endDateTime := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

	event := &entity.Event{
		Organizer:   req.Organizer,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		StartTime:   startDateTime,
		EndTime:     endDateTime,
		Location:    req.Location,
	}

	if err := s.eventRepo.CreateEvent(event); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewEventResponse(event), nil, nil
}

func (s *eventService) GetEventByID(ctx context.Context, id uuid.UUID) (*response.EventResponse, error) {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrEventNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return response.NewEventResponse(event), nil
}

func (s *eventService) GetEvents(ctx context.Context, req *request.GetEventsRequest) (*response.EventListResponse, map[string]string, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	events, total, err := s.eventRepo.GetEvents(req.Page, req.Limit, req.Search, req.OrderBy)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if len(events) == 0 {
		return &response.EventListResponse{
			Events: []*response.EventResponse{},
			Pagination: &utils_response.Pagination{
				Page:       req.Page,
				Limit:      req.Limit,
				TotalPages: 0,
				Total:      total,
			},
		}, nil, nil
	}

	var eventResponses []*response.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, response.NewEventResponse(event))
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &response.EventListResponse{
		Events: eventResponses,
		Pagination: &utils_response.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
			Total:      total,
		},
	}, nil, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, id uuid.UUID, req *request.UpdateEventRequest) (*response.EventResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	title, _ := s.eventRepo.GetEventByTitle(req.Title)
	if title != nil {
		return nil, nil, errs.ErrEventTitleAlreadyExists
	}

	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errs.ErrEventNotFound
		}
		return nil, nil, errs.ErrInternalServerError
	}

	if req.Organizer != "" {
		event.Organizer = req.Organizer
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.StartDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, map[string]string{"start_date": "Invalid start date format. Use YYYY-MM-DD"}, nil
		}
		event.StartDate = parsedDate
		if event.StartDate.After(event.EndDate) {
			return nil, map[string]string{"start_date": "Start date must be before end date"}, nil
		}
	}
	if req.EndDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, map[string]string{"end_date": "Invalid end date format. Use YYYY-MM-DD"}, nil
		}
		event.EndDate = parsedDate
		if event.EndDate.Before(event.StartDate) {
			return nil, map[string]string{"end_date": "End date must be after start date"}, nil
		}
	}
	if req.StartTime != "" {
		parsedTime, err := time.Parse("15:04", req.StartTime)
		if err != nil {
			return nil, map[string]string{"start_time": "Invalid start time format. Use HH:MM"}, nil
		}
		baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
		event.StartTime = time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	}
	if req.EndTime != "" {
		parsedTime, err := time.Parse("15:04", req.EndTime)
		if err != nil {
			return nil, map[string]string{"end_time": "Invalid end time format. Use HH:MM"}, nil
		}
		baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
		event.EndTime = time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	}

	if req.Location != "" {
		event.Location = req.Location
	}

	if event.Organizer == "" && event.Title == "" && event.Description == "" && event.StartDate.Equal(event.EndDate) && event.StartTime.Equal(event.EndTime) && event.Location == "" {
		return nil, nil, errs.ErrAtleastOneField
	}

	err = s.eventRepo.UpdateEvent(event)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewEventResponse(event), nil, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrEventNotFound
		}
		return errs.ErrInternalServerError
	}

	if err := s.eventRepo.DeleteEvent(id); err != nil {
		return errs.ErrInternalServerError
	}

	return nil
}
