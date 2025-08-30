package controller

import (
	"net/http"
	"ticert/dto/request"
	"ticert/service"
	"ticert/utils/errs"
	"ticert/utils/response"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (h *ReportController) GenerateSummaryReport(ctx *gin.Context) {
	var req request.GenerateReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	reportResponse, validationErrors, err := h.reportService.GenerateSummaryReport(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Summary report generated successfully", reportResponse, nil)
}

func (h *ReportController) GetReportList(ctx *gin.Context) {
	var req request.ReportFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	reports, validationErrors, err := h.reportService.GetReportList(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Reports retrieved successfully", reports, nil)
}
