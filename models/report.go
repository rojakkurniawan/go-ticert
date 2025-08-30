package models

import (
	"github.com/google/uuid"
)

type EventReportData struct {
	EventID      uuid.UUID `json:"event_id"`
	TotalSold    int64     `json:"total_sold"`
	TotalRevenue float64   `json:"total_revenue"`
}

type EventTicketStats struct {
	TotalOrders  int64   `json:"total_orders"`
	TotalTickets int64   `json:"total_tickets"`
	TotalRevenue float64 `json:"total_revenue"`
}

type CategoryTicketStats struct {
	TotalOrders  int64   `json:"total_orders"`
	TotalTickets int64   `json:"total_tickets"`
	TotalRevenue float64 `json:"total_revenue"`
}

type SummaryReportData struct {
	TotalTicketsSold int64   `json:"total_tickets_sold"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalEvents      int64   `json:"total_events"`
	TotalCategories  int64   `json:"total_categories"`
}
