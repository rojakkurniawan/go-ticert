package controller

import (
	"net/http"
	"ticert/dto/request"
	"ticert/service"
	"ticert/utils/errs"
	"ticert/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) *EventController {
	return &EventController{eventService: eventService}
}

func (h *EventController) CreateEvent(ctx *gin.Context) {
	var req request.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	eventResponse, validationErrors, err := h.eventService.CreateEvent(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusCreated, "Event created successfully", eventResponse, nil)
}

func (h *EventController) GetEventByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	eventResponse, err := h.eventService.GetEventByID(ctx, id)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}
	response.BuildSuccessResponse(ctx, http.StatusOK, "Event fetched successfully", eventResponse, nil)
}

func (h *EventController) GetEvents(ctx *gin.Context) {
	var req request.GetEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	events, validationErrors, err := h.eventService.GetEvents(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}
	response.BuildSuccessResponse(ctx, http.StatusOK, "Events fetched successfully", events, nil)
}

func (h *EventController) UpdateEvent(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	var req request.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	eventResponse, validationErrors, err := h.eventService.UpdateEvent(ctx, id, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Event updated successfully", eventResponse, nil)
}

func (h *EventController) DeleteEvent(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	if err := h.eventService.DeleteEvent(ctx, id); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Event deleted successfully", nil, nil)
}
