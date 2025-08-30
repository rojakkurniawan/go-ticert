package response

import (
	"ticert/utils/response"
	"time"

	"github.com/google/uuid"
)

type SummaryReportResponse struct {
	TotalTicketsSold int64                    `json:"total_tickets_sold"`
	TotalRevenue     float64                  `json:"total_revenue"`
	TotalEvents      int64                    `json:"total_events"`
	TotalCategories  int64                    `json:"total_categories"`
	Period           string                   `json:"period"`
	GeneratedAt      time.Time                `json:"generated_at"`
	EventID          *uuid.UUID               `json:"event_id,omitempty"`
	EventTitle       *string                  `json:"event_title,omitempty"`
	EventDate        string                   `json:"event_date,omitempty"`
	Categories       []CategoryReportResponse `json:"categories,omitempty"`
}

type EventReportResponse struct {
	EventID          uuid.UUID                `json:"event_id"`
	EventTitle       string                   `json:"event_title"`
	EventDate        string                   `json:"event_date"`
	TotalTicketsSold int64                    `json:"total_tickets_sold"`
	TotalRevenue     float64                  `json:"total_revenue"`
	Categories       []CategoryReportResponse `json:"categories"`
}

type CategoryReportResponse struct {
	CategoryID     uuid.UUID `json:"category_id"`
	CategoryName   string    `json:"category_name"`
	TicketsSold    int64     `json:"tickets_sold"`
	Revenue        float64   `json:"revenue"`
	RemainingStock int       `json:"remaining_stock"`
	Status         string    `json:"status"`
}

type ReportListResponse struct {
	Reports    []*EventReportResponse `json:"reports"`
	Pagination *response.Pagination   `json:"pagination"`
}
