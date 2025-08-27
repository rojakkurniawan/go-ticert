package service

import (
	"context"
	"ticert/dto/request"
	"ticert/dto/response"
	"ticert/entity"
	"ticert/repository"
	"ticert/utils/auth"
	"ticert/utils/errs"
	"ticert/utils/jwt"
	"ticert/utils/validator"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, map[string]string, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, map[string]string, error)
	RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.AuthResponse, map[string]string, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userCtx auth.ContextKey, req *request.UpdateUserRequest) (*response.UserResponse, map[string]string, error)
	UpdatePassword(ctx context.Context, userCtx auth.ContextKey, req *request.UpdatePasswordRequest) (*response.UserResponse, map[string]string, error)
	UpdateEmail(ctx context.Context, userCtx auth.ContextKey, req *request.UpdateEmailRequest) (*response.UserResponse, map[string]string, error)
	DeleteUser(ctx context.Context, userCtx auth.ContextKey) error
}

type userService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewUserService(userRepo repository.UserRepository, authRepo repository.AuthRepository) UserService {
	return &userService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *userService) Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, nil, errs.ErrEmailAlreadyExists
	}

	user := &entity.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	accessExpiry, _ := jwt.GetAccessTokenExpiry()
	refreshExpiry, _ := jwt.GetRefreshTokenExpiry()

	if err := s.authRepo.StoreAccessToken(user.ID, accessToken, accessExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if err := s.authRepo.StoreRefreshToken(user.ID, refreshToken, refreshExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewAuthResponse(accessToken, refreshToken, int64(accessExpiry.Seconds()), user), nil, nil
}

func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, nil, errs.ErrAuthInvalidCredentials
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, nil, errs.ErrAuthInvalidCredentials
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	s.authRepo.RevokeAllUserTokens(user.ID)

	accessExpiry, _ := jwt.GetAccessTokenExpiry()
	refreshExpiry, _ := jwt.GetRefreshTokenExpiry()

	if err := s.authRepo.StoreAccessToken(user.ID, accessToken, accessExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if err := s.authRepo.StoreRefreshToken(user.ID, refreshToken, refreshExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewAuthResponse(accessToken, refreshToken, int64(accessExpiry.Seconds()), user), nil, nil
}

func (s *userService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.AuthResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	claims, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, nil, errs.ErrInvalidRefreshToken
	}

	valid, err := s.authRepo.ValidateRefreshToken(claims.UserID, req.RefreshToken)
	if err != nil || !valid {
		return nil, nil, errs.ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, nil, errs.ErrUserNotFound
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	accessExpiry, _ := jwt.GetAccessTokenExpiry()
	refreshExpiry, _ := jwt.GetRefreshTokenExpiry()

	if err := s.authRepo.StoreAccessToken(user.ID, accessToken, accessExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if err := s.authRepo.StoreRefreshToken(user.ID, refreshToken, refreshExpiry); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewAuthResponse(accessToken, refreshToken, int64(accessExpiry.Seconds()), user), nil, nil
}

func (s *userService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.authRepo.RevokeAllUserTokens(userID)
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}
	return response.NewUserResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, userCtx auth.ContextKey, req *request.UpdateUserRequest) (*response.UserResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	user, err := s.userRepo.GetUserByID(userCtx.UserID)
	if err != nil {
		return nil, nil, errs.ErrUserNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if userCtx.Role == "admin" {
		if req.Role != nil {
			user.Role = *req.Role
		}
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewUserResponse(updatedUser), nil, nil
}

func (s *userService) UpdatePassword(ctx context.Context, userCtx auth.ContextKey, req *request.UpdatePasswordRequest) (*response.UserResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	user, err := s.userRepo.GetUserByID(userCtx.UserID)
	if err != nil {
		return nil, nil, errs.ErrUserNotFound
	}

	if err := user.CheckPassword(req.OldPassword); err != nil {
		return nil, nil, errs.ErrAuthInvalidCredentials
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	if err := s.authRepo.RevokeAllUserTokens(user.ID); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	return response.NewUserResponse(updatedUser), nil, nil

}

func (s *userService) UpdateEmail(ctx context.Context, userCtx auth.ContextKey, req *request.UpdateEmailRequest) (*response.UserResponse, map[string]string, error) {
	validationErrors := validator.HandleValidationErrors(req)
	if validationErrors != nil {
		return nil, validationErrors, nil
	}

	user, err := s.userRepo.GetUserByID(userCtx.UserID)
	if err != nil {
		return nil, nil, errs.ErrUserNotFound
	}

	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, nil, errs.ErrEmailAlreadyExists
	}

	user.Email = req.Email

	if err := s.authRepo.RevokeAllUserTokens(user.ID); err != nil {
		return nil, nil, errs.ErrInternalServerError
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, nil, err
	}

	return response.NewUserResponse(updatedUser), nil, nil
}

func (s *userService) DeleteUser(ctx context.Context, userCtx auth.ContextKey) error {
	user, err := s.userRepo.GetUserByID(userCtx.UserID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if err := s.userRepo.DeleteUser(user.ID); err != nil {
		return errs.ErrInternalServerError
	}

	return nil
}
