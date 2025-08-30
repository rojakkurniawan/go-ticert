package response

import (
	"ticert/entity"
	"ticert/utils/response"
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID         uuid.UUID              `json:"id"`
	InvoiceID  string                 `json:"invoice_id"`
	Status     string                 `json:"status"`
	Quantity   int                    `json:"quantity"`
	TotalPrice float64                `json:"total_price"`
	Tickets    []*OrderDetailResponse `json:"tickets,omitempty"`
	Event      *EventResponse         `json:"event"`
	Category   *CategoryResponse      `json:"category"`
	User       *UserResponse          `json:"user"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type OrderDetailResponse struct {
	ID             uuid.UUID `json:"id"`
	TicketCode     string    `json:"ticket_code"`
	FullName       string    `json:"full_name"`
	IdentityNumber string    `json:"identity_number"`
	Redeemed       bool      `json:"redeemed"`
}

type OrderListResponse struct {
	Orders     []*OrderResponse     `json:"orders"`
	Pagination *response.Pagination `json:"pagination"`
}

func NewOrderResponse(order *entity.Order) *OrderResponse {
	var tickets []*OrderDetailResponse
	if order.Status != "pending" && order.Status != "cancelled" {
		tickets = NewOrderDetailListResponse(order.OrderDetails)
	}
	return &OrderResponse{
		ID:         order.ID,
		InvoiceID:  order.InvoiceID,
		Status:     order.Status,
		Quantity:   order.Quantity,
		TotalPrice: order.TotalPrice,
		Tickets:    tickets,
		Event:      NewEventResponse(order.Category.Event),
		Category: &CategoryResponse{
			ID:        order.Category.ID,
			Name:      order.Category.Name,
			Price:     order.Category.Price,
			EventDate: order.Category.EventDate.Format("02 Jan 2006"),
		},
		User:      NewUserResponse(order.User),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func NewOrderDetailListResponse(orderDetails []*entity.OrderDetail) []*OrderDetailResponse {
	orderDetailsResponse := make([]*OrderDetailResponse, len(orderDetails))
	for i, orderDetail := range orderDetails {
		orderDetailsResponse[i] = NewOrderDetailResponse(orderDetail)
	}
	return orderDetailsResponse
}

func NewOrderDetailResponse(orderDetail *entity.OrderDetail) *OrderDetailResponse {
	return &OrderDetailResponse{
		ID:             orderDetail.ID,
		TicketCode:     orderDetail.TicketCode,
		FullName:       orderDetail.FullName,
		IdentityNumber: orderDetail.IdentityNumber,
		Redeemed:       orderDetail.Redeemed,
	}
}
