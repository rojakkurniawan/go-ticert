package service

import (
	"context"
	"ticert/dto/request"
	"ticert/dto/response"
	"ticert/repository"
	"ticert/utils/errs"
	utils_response "ticert/utils/response"
	"ticert/utils/validator"
	"time"

	"github.com/google/uuid"
)

type ReportService interface {
	GenerateSummaryReport(ctx context.Context, req *request.GenerateReportRequest) (*response.SummaryReportResponse, map[string]string, error)
	GetReportList(ctx context.Context, req *request.ReportFilterRequest) (*response.ReportListResponse, map[string]string, error)
}

type reportService struct {
	reportRepo   repository.ReportRepository
	eventRepo    repository.EventRepository
	categoryRepo repository.CategoryRepository
}

func NewReportService(reportRepo repository.ReportRepository, eventRepo repository.EventRepository, categoryRepo repository.CategoryRepository) ReportService {
	return &reportService{
		reportRepo:   reportRepo,
		eventRepo:    eventRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *reportService) GenerateSummaryReport(ctx context.Context, req *request.GenerateReportRequest) (*response.SummaryReportResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	var startDate, endDate *time.Time
	var err error

	if req.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, map[string]string{"start_date": "Invalid start date format. Use YYYY-MM-DD"}, nil
		}
		startDate = &parsedStartDate
	}

	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, map[string]string{"end_date": "Invalid end date format. Use YYYY-MM-DD"}, nil
		}
		endDate = &parsedEndDate
	}

	if req.EventID != "" {
		eventID, err := uuid.Parse(req.EventID)
		if err != nil {
			return nil, map[string]string{"event_id": "Invalid event ID format"}, nil
		}

		event, err := s.eventRepo.GetEventByID(eventID)
		if err != nil {
			return nil, nil, errs.ErrEventNotFound
		}

		eventStats, err := s.reportRepo.GetEventTicketStats(eventID)
		if err != nil {
			return nil, nil, errs.ErrInternalServerError
		}

		categories, err := s.categoryRepo.GetCategories(eventID)
		if err != nil {
			return nil, nil, errs.ErrInternalServerError
		}

		var categoryReports []response.CategoryReportResponse
		for _, category := range categories {
			categoryStats, err := s.reportRepo.GetCategoryTicketStats(category.ID)
			if err != nil {
				continue
			}

			categoryReports = append(categoryReports, response.CategoryReportResponse{
				CategoryID:     category.ID,
				CategoryName:   category.Name,
				TicketsSold:    categoryStats.TotalTickets,
				Revenue:        categoryStats.TotalRevenue,
				RemainingStock: category.Quantity,
				Status:         category.Status,
			})
		}

		period := "All Time"
		if startDate != nil && endDate != nil {
			period = startDate.Format("02 Jan 2006") + " - " + endDate.Format("02 Jan 2006")
		} else if startDate != nil {
			period = "From " + startDate.Format("02 Jan 2006")
		} else if endDate != nil {
			period = "Until " + endDate.Format("02 Jan 2006")
		}

		return &response.SummaryReportResponse{
			TotalTicketsSold: eventStats.TotalTickets,
			TotalRevenue:     eventStats.TotalRevenue,
			TotalEvents:      1, // Hanya 1 event
			TotalCategories:  int64(len(categories)),
			Period:           period,
			GeneratedAt:      time.Now(),
			EventID:          &event.ID,
			EventTitle:       &event.Title,
			EventDate:        getEventDate(event.StartDate, event.EndDate),
			Categories:       categoryReports,
		}, nil, nil
	}

	totalTickets, err := s.reportRepo.GetTotalTicketsSold(startDate, endDate)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	totalRevenue, err := s.reportRepo.GetTotalRevenue(startDate, endDate)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	totalEvents, err := s.reportRepo.GetTotalEvents()
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	totalCategories, err := s.reportRepo.GetTotalCategories()
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	period := "All Time"
	if startDate != nil && endDate != nil {
		period = startDate.Format("02 Jan 2006") + " - " + endDate.Format("02 Jan 2006")
	} else if startDate != nil {
		period = "From " + startDate.Format("02 Jan 2006")
	} else if endDate != nil {
		period = "Until " + endDate.Format("02 Jan 2006")
	}

	return &response.SummaryReportResponse{
		TotalTicketsSold: totalTickets,
		TotalRevenue:     totalRevenue,
		TotalEvents:      totalEvents,
		TotalCategories:  totalCategories,
		Period:           period,
		GeneratedAt:      time.Now(),
	}, nil, nil
}

func (s *reportService) GetReportList(ctx context.Context, req *request.ReportFilterRequest) (*response.ReportListResponse, map[string]string, error) {
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

	var eventID *uuid.UUID
	var startDate, endDate *time.Time
	var err error

	if req.EventID != "" {
		parsedEventID, err := uuid.Parse(req.EventID)
		if err != nil {
			return nil, map[string]string{"event_id": "Invalid event ID format"}, nil
		}
		eventID = &parsedEventID
	}

	if req.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, map[string]string{"start_date": "Invalid start date format. Use YYYY-MM-DD"}, nil
		}
		startDate = &parsedStartDate
	}

	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, map[string]string{"end_date": "Invalid end date format. Use YYYY-MM-DD"}, nil
		}
		endDate = &parsedEndDate
	}

	eventReports, err := s.reportRepo.GetEventReports(eventID, startDate, endDate)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	var eventReportResponses []*response.EventReportResponse
	for _, report := range eventReports {
		event, err := s.eventRepo.GetEventByID(report.EventID)
		if err != nil {
			continue
		}

		categories, err := s.categoryRepo.GetCategories(report.EventID)
		if err != nil {
			continue
		}

		var categoryReports []response.CategoryReportResponse
		for _, category := range categories {
			categoryStats, err := s.reportRepo.GetCategoryTicketStats(category.ID)
			if err != nil {
				continue
			}

			categoryReports = append(categoryReports, response.CategoryReportResponse{
				CategoryID:     category.ID,
				CategoryName:   category.Name,
				TicketsSold:    categoryStats.TotalTickets,
				Revenue:        categoryStats.TotalRevenue,
				RemainingStock: category.Quantity,
				Status:         category.Status,
			})
		}

		eventReportResponses = append(eventReportResponses, &response.EventReportResponse{
			EventID:          event.ID,
			EventTitle:       event.Title,
			EventDate:        getEventDate(event.StartDate, event.EndDate),
			TotalTicketsSold: report.TotalSold,
			TotalRevenue:     report.TotalRevenue,
			Categories:       categoryReports,
		})
	}

	total := int64(len(eventReportResponses))
	start := (req.Page - 1) * req.Limit
	end := start + req.Limit
	if end > int(total) {
		end = int(total)
	}
	if start > int(total) {
		start = int(total)
	}

	var paginatedReports []*response.EventReportResponse
	if start < int(total) {
		paginatedReports = eventReportResponses[start:end]
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &response.ReportListResponse{
		Reports: paginatedReports,
		Pagination: &utils_response.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
			Total:      total,
		},
	}, nil, nil
}

func getEventDate(startDate, endDate time.Time) string {
	if startDate.Month() == endDate.Month() && startDate.Year() == endDate.Year() {
		return startDate.Format("02") + " - " + endDate.Format("02 Jan 2006")
	}
	return startDate.Format("02 Jan 2006") + " - " + endDate.Format("02 Jan 2006")
}
