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

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (h *CategoryController) CreateCategory(ctx *gin.Context) {
	var req request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	categoryResponse, validationErrors, err := h.categoryService.CreateCategory(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusCreated, "Category created successfully", categoryResponse, nil)
}

func (h *CategoryController) GetCategoryByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	categoryResponse, err := h.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Category retrieved successfully", categoryResponse, nil)
}

func (h *CategoryController) GetCategories(ctx *gin.Context) {
	eventIDParam := ctx.Param("event_id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	categories, err := h.categoryService.GetCategories(ctx, eventID)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Categories retrieved successfully", categories, nil)
}

func (h *CategoryController) UpdateCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	var req request.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	categoryResponse, validationErrors, err := h.categoryService.UpdateCategory(ctx, id, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Category updated successfully", categoryResponse, nil)
}

func (h *CategoryController) DeleteCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	if err := h.categoryService.DeleteCategory(ctx, id); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Category deleted successfully", nil, nil)
}
