package request

import (
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	EventID   uuid.UUID `json:"event_id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required,max=255"`
	Price     float64   `json:"price" validate:"required,min=0"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	EventDate string    `json:"event_date" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name      string  `json:"name" validate:"omitempty,max=255"`
	Price     float64 `json:"price" validate:"omitempty,min=0"`
	Quantity  int     `json:"quantity" validate:"omitempty,min=1"`
	EventDate string  `json:"event_date" validate:"omitempty"`
}
