package response

import (
	"ticert/entity"
	"ticert/utils/response"
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID          uuid.UUID           `json:"id"`
	Organizer   string              `json:"organizer"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	EventDate   string              `json:"event_date"`
	RangeTime   string              `json:"range_time"`
	Location    string              `json:"location"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Categories  []*CategoryResponse `json:"categories,omitempty"`
}

type EventListResponse struct {
	Events     []*EventResponse     `json:"events"`
	Pagination *response.Pagination `json:"pagination"`
}

func NewEventResponse(event *entity.Event) *EventResponse {
	return &EventResponse{
		ID:          event.ID,
		Organizer:   event.Organizer,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   getEventDate(event.StartDate, event.EndDate),
		RangeTime:   event.StartTime.Format("15:04") + " - " + event.EndTime.Format("15:04"),
		Location:    event.Location,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		Categories:  NewCategoryListResponse(event.Categories),
	}
}

func getEventDate(startDate, endDate time.Time) string {
	if startDate.Month() == endDate.Month() && startDate.Year() == endDate.Year() {
		return startDate.Format("02") + " - " + endDate.Format("02 Jan 2006")
	}
	return startDate.Format("02 Jan 2006") + " - " + endDate.Format("02 Jan 2006")
}
