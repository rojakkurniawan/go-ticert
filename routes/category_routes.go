package routes

import (
	"ticert/controller"
	"ticert/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(r *gin.Engine, categoryController *controller.CategoryController) {
	protected := r.Group("/api/v1/categories")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/", middleware.RoleMiddleware("admin"), categoryController.CreateCategory)
		protected.GET("/:id", categoryController.GetCategoryByID)
		protected.GET("/event/:event_id", categoryController.GetCategories)
		protected.PATCH("/:id", middleware.RoleMiddleware("admin"), categoryController.UpdateCategory)
		protected.DELETE("/:id", middleware.RoleMiddleware("admin"), categoryController.DeleteCategory)
	}
}
