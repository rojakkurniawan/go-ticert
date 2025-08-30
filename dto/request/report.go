package request

type GenerateReportRequest struct {
	EventID    string `form:"event_id" validate:"omitempty"`
	CategoryID string `form:"category_id" validate:"omitempty"`
	StartDate  string `form:"start_date" validate:"omitempty"`
	EndDate    string `form:"end_date" validate:"omitempty"`
}

type ReportFilterRequest struct {
	EventID    string `form:"event_id" validate:"omitempty"`
	CategoryID string `form:"category_id" validate:"omitempty"`
	StartDate  string `form:"start_date" validate:"omitempty"`
	EndDate    string `form:"end_date" validate:"omitempty"`
	Page       int    `form:"page" validate:"omitempty,min=1"`
	Limit      int    `form:"limit" validate:"omitempty,min=1,max=100"`
}
