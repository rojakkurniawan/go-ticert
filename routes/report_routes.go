package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(r *gin.Engine, reportController *controller.ReportController) {
	protected := r.Group("/api/v1/reports")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.RoleMiddleware("admin"))

	{
		protected.GET("/summary", reportController.GenerateSummaryReport)
		protected.GET("/", reportController.GetReportList)
	}
}
