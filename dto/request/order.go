package request

import "github.com/google/uuid"

type OrderRequest struct {
	CategoryID    uuid.UUID            `json:"category_id" validate:"required,uuid"`
	Quantity      int                  `json:"quantity" validate:"required,min=1,max=10"`
	SameAsOrderer bool                 `json:"same_as_orderer" validate:"omitempty"`
	OrderDetails  []OrderDetailRequest `json:"order_details" validate:"required,min=1,max=10"`
}

type OrderDetailRequest struct {
	FullName       string `json:"full_name" validate:"required,min=3,max=255"`
	IdentityNumber string `json:"identity_number" validate:"required,min=5,max=20"`
}

type GetOrdersRequest struct {
	Page  int `form:"page" validate:"omitempty,min=1"`
	Limit int `form:"limit" validate:"omitempty,min=1"`
}

type GetOrdersRequestAdmin struct {
	Page    int    `form:"page" validate:"omitempty"`
	Limit   int    `form:"limit" validate:"omitempty"`
	Status  string `form:"status" validate:"omitempty,oneof=pending paid cancelled"`
	Search  string `form:"search" validate:"omitempty,max=255"`
	OrderBy string `form:"order_by" validate:"omitempty,oneof=asc desc quantity_asc quantity_desc total_price_asc total_price_desc"`
}
