package response

import (
	"ticert/entity"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	EventDate string    `json:"event_date"`
	Price     float64   `json:"price,omitempty"`
	Quantity  int       `json:"quantity,omitempty"`
	Status    string    `json:"status,omitempty"`
}

func NewCategoryResponse(ticket *entity.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:        ticket.ID,
		Name:      ticket.Name,
		EventDate: ticket.EventDate.Format("02 Jan 2006"),
		Price:     ticket.Price,
		Quantity:  ticket.Quantity,
		Status:    ticket.Status,
	}
}

func NewCategoryListResponse(categories []entity.Category) []*CategoryResponse {
	var responses []*CategoryResponse
	for _, category := range categories {
		responses = append(responses, NewCategoryResponse(&category))
	}
	return responses
}
