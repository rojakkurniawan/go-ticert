package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userController *controller.UserController) {
	protected := r.Group("/api/v1/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", userController.GetProfile)
		protected.PATCH("/profile", userController.UpdateProfile)
		protected.PATCH("/password", userController.UpdatePassword)
		protected.PATCH("/email", userController.UpdateEmail)
		protected.DELETE("/", userController.DeleteUser)
	}
}
