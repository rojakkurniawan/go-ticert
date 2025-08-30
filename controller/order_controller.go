package controller

import (
	"net/http"
	"ticert/dto/request"
	"ticert/service"
	"ticert/utils/auth"
	"ticert/utils/errs"
	"ticert/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	orderService service.OrderService
}

func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (h *OrderController) CreateOrder(ctx *gin.Context) {
	var req request.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	orderResponse, validationErrors, err := h.orderService.CreateOrder(ctx, &req, userCtx.UserID)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusCreated, "Order created successfully", orderResponse, nil)
}

func (h *OrderController) GetOrders(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	var req request.GetOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	orders, validationErrors, err := h.orderService.GetOrders(ctx, userCtx.UserID, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Orders fetched successfully", orders, nil)
}

func (h *OrderController) GetOrdersAdmin(ctx *gin.Context) {
	var req request.GetOrdersRequestAdmin
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	orders, validationErrors, err := h.orderService.GetOrdersAdmin(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Orders fetched successfully", orders, nil)
}

func (h *OrderController) GetOrderById(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	order, err := h.orderService.GetOrderById(ctx, orderID, &userCtx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Order fetched successfully", order, nil)
}

func (h *OrderController) CancelOrder(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}
	if err := h.orderService.CancelOrder(ctx, orderID, userCtx.UserID); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Order cancelled successfully", nil, nil)
}

func (h *OrderController) VerifyOrderStatus(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	if err := h.orderService.VerifyOrderStatus(ctx, orderID); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Order verified successfully", nil, nil)
}

func (h *OrderController) VerifyTicket(ctx *gin.Context) {
	ticketCode := ctx.Param("ticket_code")

	if err := h.orderService.VerifyTicket(ctx, ticketCode); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Ticket redeemed successfully", nil, nil)
}
