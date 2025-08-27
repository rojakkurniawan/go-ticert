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

	userService := service.NewUserService(userRepo, authRepo)

	userController := controller.NewUserController(userService)

	SetupAuthRoutes(r, userController)
	SetupUserRoutes(r, userController)
}
