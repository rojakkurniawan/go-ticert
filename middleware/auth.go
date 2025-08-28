package middleware

import (
	"strings"
	"ticert/repository"
	"ticert/utils/errs"
	"ticert/utils/jwt"
	"ticert/utils/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.BuildErrorResponse(c, errs.ErrLoginRequired)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == authHeader {
			response.BuildErrorResponse(c, errs.ErrInvalidAccessToken)
			c.Abort()
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			response.BuildErrorResponse(c, errs.ErrSessionExpired)
			c.Abort()
			return
		}

		authRepo := repository.NewAuthRepository()
		valid, err := authRepo.ValidateAccessToken(claims.UserID, tokenString)
		if err != nil || !valid {
			response.BuildErrorResponse(c, errs.ErrSessionExpired)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.BuildErrorResponse(c, errs.ErrAccessDenied)
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		response.BuildErrorResponse(c, errs.ErrInsufficientPermissions)
		c.Abort()
	}
}
