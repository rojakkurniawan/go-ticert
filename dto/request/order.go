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
