package service

import (
	"context"
	"errors"
	"ticert/dto/request"
	"ticert/dto/response"
	"ticert/entity"
	"ticert/repository"
	"ticert/utils/errs"
	utils_response "ticert/utils/response"
	"ticert/utils/validator"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *request.OrderRequest, userID uuid.UUID) (*response.OrderResponse, map[string]string, error)
	GetOrders(ctx context.Context, userID uuid.UUID, req *request.GetOrdersRequest) (*response.OrderListResponse, map[string]string, error)
	GetOrderById(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*response.OrderResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error
	VerifyOrderStatus(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error
	VerifyTicket(ctx context.Context, ticketCode string) error
}

type orderService struct {
	orderRepository    repository.OrderRepository
	userRepository     repository.UserRepository
	categoryRepository repository.CategoryRepository
}

func NewOrderService(orderRepository repository.OrderRepository, userRepository repository.UserRepository, categoryRepository repository.CategoryRepository) OrderService {
	return &orderService{orderRepository: orderRepository, userRepository: userRepository, categoryRepository: categoryRepository}
}

func (s *orderService) CreateOrder(ctx context.Context, req *request.OrderRequest, userID uuid.UUID) (*response.OrderResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, nil, err
	}

	category, err := s.categoryRepository.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, nil, err
	}

	stockAvailable, err := s.categoryRepository.CheckStock(req.CategoryID, req.Quantity)
	if err != nil {
		return nil, nil, err
	}

	if !stockAvailable {
		return nil, nil, errs.ErrStockNotAvailable
	}

	order := &entity.Order{
		UserID:     user.ID,
		CategoryID: req.CategoryID,
		Quantity:   req.Quantity,
		TotalPrice: float64(req.Quantity) * category.Price,
		Status:     "pending",
		InvoiceID:  uuid.New().String()[26:],
	}

	var orderDetails []*entity.OrderDetail

	if req.SameAsOrderer {
		if len(req.OrderDetails) > 0 {
			firstDetail := req.OrderDetails[0]
			for i := 0; i < req.Quantity; i++ {
				orderDetail := &entity.OrderDetail{
					TicketCode:     uuid.New().String()[26:],
					FullName:       firstDetail.FullName,
					IdentityNumber: firstDetail.IdentityNumber,
				}
				orderDetails = append(orderDetails, orderDetail)
			}
		}
	} else {
		if len(req.OrderDetails) != req.Quantity {
			return nil, nil, errs.ErrQuantityNotMatch
		}

		orderDetails = make([]*entity.OrderDetail, len(req.OrderDetails))
		for i, detail := range req.OrderDetails {
			orderDetails[i] = &entity.OrderDetail{
				TicketCode:     uuid.New().String()[26:],
				FullName:       detail.FullName,
				IdentityNumber: detail.IdentityNumber,
			}
		}
	}

	if err := s.orderRepository.CreateOrder(order, orderDetails, req.CategoryID, req.Quantity); err != nil {
		return nil, nil, err
	}

	return response.NewOrderResponse(order), nil, nil
}

func (s *orderService) GetOrders(ctx context.Context, userID uuid.UUID, req *request.GetOrdersRequest) (*response.OrderListResponse, map[string]string, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	orders, total, err := s.orderRepository.GetOrders(req.Page, req.Limit, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(orders) == 0 {
		return &response.OrderListResponse{
			Orders: []*response.OrderResponse{},
			Pagination: &utils_response.Pagination{
				Page:       req.Page,
				Limit:      req.Limit,
				TotalPages: 0,
				Total:      total,
			},
		}, nil, nil
	}

	var orderResponses []*response.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, response.NewOrderResponse(order))
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &response.OrderListResponse{
		Orders: orderResponses,
		Pagination: &utils_response.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
			Total:      total,
		},
	}, nil, nil
}

func (s *orderService) GetOrderById(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*response.OrderResponse, error) {
	order, err := s.orderRepository.GetOrderById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, errs.ErrInternalServerError
	}

	if order.UserID != userID {
		return nil, errs.ErrUnauthorized
	}

	return response.NewOrderResponse(order), nil
}

func (s *orderService) CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepository.GetOrderById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrOrderNotFound
		}
		return errs.ErrInternalServerError
	}

	if order.UserID != userID {
		return errs.ErrUnauthorized
	}

	if order.Status == "cancelled" {
		return errs.ErrOrderAlreadyCancelled
	}

	if order.Status == "paid" {
		return errs.ErrOrderAlreadyPaid
	}

	if err := s.orderRepository.CancelOrder(orderID); err != nil {
		return err
	}

	return nil
}

func (s *orderService) VerifyOrderStatus(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepository.GetOrderById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrOrderNotFound
		}
		return errs.ErrInternalServerError
	}

	if order.UserID != userID {
		return errs.ErrUnauthorized
	}

	if order.Status == "cancelled" {
		return errs.ErrOrderAlreadyCancelled
	}

	if order.Status == "paid" {
		return errs.ErrOrderAlreadyPaid
	}

	if err := s.orderRepository.VerifyOrderStatus(orderID); err != nil {
		return err
	}

	return nil
}

func (s *orderService) VerifyTicket(ctx context.Context, ticketCode string) error {
	orderDetail, err := s.orderRepository.GetOrderDetailByTicketCode(ticketCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrTicketNotFound
		}
		return errs.ErrInternalServerError
	}

	if orderDetail.Order.Status == "pending" {
		return errs.ErrOrderNotPaid
	}

	if orderDetail.Order.Status == "cancelled" {
		return errs.ErrOrderAlreadyCancelled
	}

	if orderDetail.Redeemed {
		return errs.ErrTicketAlreadyRedeemed
	}

	if err := s.orderRepository.VerifyTicket(orderDetail.ID); err != nil {
		return err
	}

	return nil
}
