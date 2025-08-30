package routes

import (
	"ticert/controller"
	"ticert/repository"
	"ticert/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	SetupMiddleware(r)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository()
	eventRepo := repository.NewEventRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reportRepo := repository.NewReportRepository(db)

	userService := service.NewUserService(userRepo, authRepo)
	eventService := service.NewEventService(eventRepo)
	categoryService := service.NewCategoryService(categoryRepo, eventRepo)
	orderService := service.NewOrderService(orderRepo, userRepo, categoryRepo)
	reportService := service.NewReportService(reportRepo, eventRepo, categoryRepo)

	userController := controller.NewUserController(userService)
	eventController := controller.NewEventController(eventService)
	categoryController := controller.NewCategoryController(categoryService)
	orderController := controller.NewOrderController(orderService)
	reportController := controller.NewReportController(reportService)

	SetupAuthRoutes(r, userController)
	SetupUserRoutes(r, userController)
	SetupEventRoutes(r, eventController)
	SetupCategoryRoutes(r, categoryController)
	SetupOrderRoutes(r, orderController)
	SetupReportRoutes(r, reportController)
}
