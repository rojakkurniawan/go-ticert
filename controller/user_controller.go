package controller

import (
	"net/http"
	"ticert/dto/request"
	"ticert/service"
	"ticert/utils/auth"
	"ticert/utils/errs"
	"ticert/utils/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	authResponse, validationErrors, err := h.userService.Register(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusCreated, "User registered successfully", authResponse, nil)
}

func (h *UserController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	authResponse, validationErrors, err := h.userService.Login(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Login successful", authResponse, nil)
}

func (h *UserController) RefreshToken(ctx *gin.Context) {
	var req request.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	authResponse, validationErrors, err := h.userService.RefreshToken(ctx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Token refreshed successfully", authResponse, nil)
}

func (h *UserController) Logout(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	if err := h.userService.Logout(ctx, userCtx.UserID); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Logout successful", nil, nil)
}

func (h *UserController) GetProfile(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	userResponse, err := h.userService.GetUserByID(ctx, userCtx.UserID)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", userResponse, nil)
}

func (h *UserController) UpdateProfile(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	userResponse, validationErrors, err := h.userService.UpdateUser(ctx, userCtx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Profile updated successfully", userResponse, nil)
}

func (h *UserController) UpdatePassword(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	var req request.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	_, validationErrors, err := h.userService.UpdatePassword(ctx, userCtx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Password updated successfully", nil, nil)
}

func (h *UserController) UpdateEmail(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	var req request.UpdateEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BuildErrorResponse(ctx, errs.ErrBadRequest)
		return
	}

	userResponse, validationErrors, err := h.userService.UpdateEmail(ctx, userCtx, &req)
	if validationErrors != nil {
		response.BuildValidationErrorResponse(ctx, validationErrors)
		return
	}

	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "Email updated successfully", userResponse, nil)
}

func (h *UserController) DeleteUser(ctx *gin.Context) {
	userCtx, err := auth.GetUserContextKey(ctx)
	if err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	if err := h.userService.DeleteUser(ctx, userCtx); err != nil {
		response.BuildErrorResponse(ctx, err)
		return
	}

	response.BuildSuccessResponse(ctx, http.StatusOK, "User deleted successfully", nil, nil)
}
