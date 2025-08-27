package auth

import (
	"ticert/utils/errs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

func GetUserContextKey(ctx *gin.Context) (ContextKey, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return ContextKey{}, errs.ErrLoginRequired
	}

	role, exists := ctx.Get("role")
	if !exists {
		return ContextKey{}, errs.ErrLoginRequired
	}

	email, exists := ctx.Get("email")
	if !exists {
		return ContextKey{}, errs.ErrLoginRequired
	}

	return ContextKey{
		UserID: userID.(uuid.UUID),
		Email:  email.(string),
		Role:   role.(string),
	}, nil
}
