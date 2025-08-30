package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(r *gin.Engine, eventController *controller.EventController) {
	protected := r.Group("/api/v1/events")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/", middleware.RoleMiddleware("admin"), eventController.CreateEvent)
		protected.GET("/:id", eventController.GetEventByID)
		protected.GET("/", eventController.GetEvents)
		protected.PATCH("/:id", middleware.RoleMiddleware("admin"), eventController.UpdateEvent)
		protected.DELETE("/:id", middleware.RoleMiddleware("admin"), eventController.DeleteEvent)
	}
}
