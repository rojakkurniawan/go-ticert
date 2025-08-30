package repository

import (
	"ticert/entity"
	"ticert/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetTotalTicketsSold(startDate, endDate *time.Time) (int64, error)
	GetTotalRevenue(startDate, endDate *time.Time) (float64, error)
	GetTotalEvents() (int64, error)
	GetTotalCategories() (int64, error)
	GetEventTicketStats(eventID uuid.UUID) (*models.EventTicketStats, error)
	GetCategoryTicketStats(categoryID uuid.UUID) (*models.CategoryTicketStats, error)
	GetEventReports(eventID *uuid.UUID, startDate, endDate *time.Time) ([]*models.EventReportData, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetTotalTicketsSold(startDate, endDate *time.Time) (int64, error) {
	var total int64

	query := r.db.Model(&entity.Order{}).Where("status = ?", "paid")

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Select("COALESCE(SUM(quantity), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *reportRepository) GetTotalRevenue(startDate, endDate *time.Time) (float64, error) {
	var total float64

	query := r.db.Model(&entity.Order{}).Where("status = ?", "paid")

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Select("COALESCE(SUM(total_price), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *reportRepository) GetTotalEvents() (int64, error) {
	var total int64
	if err := r.db.Model(&entity.Event{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetTotalCategories() (int64, error) {
	var total int64
	if err := r.db.Model(&entity.Category{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetEventTicketStats(eventID uuid.UUID) (*models.EventTicketStats, error) {
	var stats models.EventTicketStats

	query := r.db.Model(&entity.Order{}).
		Select("COUNT(*) as total_orders, SUM(orders.quantity) as total_tickets, COALESCE(SUM(orders.total_price), 0) as total_revenue").
		Joins("JOIN categories ON orders.category_id = categories.id").
		Where("categories.event_id = ? AND orders.status = ?", eventID, "paid")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *reportRepository) GetCategoryTicketStats(categoryID uuid.UUID) (*models.CategoryTicketStats, error) {
	var stats models.CategoryTicketStats

	query := r.db.Model(&entity.Order{}).
		Select("COUNT(*) as total_orders, SUM(orders.quantity) as total_tickets, COALESCE(SUM(orders.total_price), 0) as total_revenue").
		Where("category_id = ? AND status = ?", categoryID, "paid")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *reportRepository) GetEventReports(eventID *uuid.UUID, startDate, endDate *time.Time) ([]*models.EventReportData, error) {
	var results []*models.EventReportData

	query := r.db.Model(&entity.Order{}).
		Select("categories.event_id, SUM(orders.quantity) as total_sold, COALESCE(SUM(orders.total_price), 0) as total_revenue").
		Joins("JOIN categories ON orders.category_id = categories.id").
		Where("orders.status = ?", "paid")

	if eventID != nil {
		query = query.Where("categories.event_id = ?", eventID)
	}
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at <= ?", endDate)
	}

	query = query.Group("categories.event_id")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
