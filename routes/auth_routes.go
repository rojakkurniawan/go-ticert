package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, userController *controller.UserController) {
	public := r.Group("/api/v1/auth")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
		public.POST("/refresh", userController.RefreshToken)
	}

	protected := r.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)
	}
}
