package request

type CreateEventRequest struct {
	Organizer   string `json:"organizer" validate:"required,max=255"`
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"required,max=255"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
	Location    string `json:"location" validate:"required,max=255"`
}

type UpdateEventRequest struct {
	Organizer   string `json:"organizer" validate:"omitempty,max=255"`
	Title       string `json:"title" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty,max=255"`
	StartDate   string `json:"start_date" validate:"omitempty"`
	EndDate     string `json:"end_date" validate:"omitempty"`
	StartTime   string `json:"start_time" validate:"omitempty"`
	EndTime     string `json:"end_time" validate:"omitempty"`
	Location    string `json:"location" validate:"omitempty,max=255"`
}

type GetEventsRequest struct {
	Page    int    `form:"page" validate:"omitempty"`
	Limit   int    `form:"limit" validate:"omitempty"`
	Search  string `form:"search" validate:"omitempty,max=255"`
	OrderBy string `form:"order_by" validate:"omitempty,oneof=asc desc date_asc date_desc time_asc time_desc"`
}
