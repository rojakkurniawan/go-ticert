package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(r *gin.Engine, orderController *controller.OrderController) {
	protected := r.Group("/api/v1/tickets")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/", orderController.CreateOrder)
		protected.GET("/", orderController.GetOrders)
		protected.GET("/:id", orderController.GetOrderById)
		protected.PATCH("/:id/cancel", orderController.CancelOrder)
		protected.GET("/admin", middleware.RoleMiddleware("admin"), orderController.GetOrdersAdmin)
		protected.PATCH("/:id/verify", middleware.RoleMiddleware("admin"), orderController.VerifyOrderStatus)
		protected.PATCH("/redeem/:ticket_code", middleware.RoleMiddleware("admin"), orderController.VerifyTicket)
	}
}
