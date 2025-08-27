package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrAccessDenied = response.ErrorModel{
		Message:    "Access denied",
		StatusCode: http.StatusForbidden,
	}

	ErrInsufficientPermissions = response.ErrorModel{
		Message:    "You don't have permission to access this resource",
		StatusCode: http.StatusForbidden,
	}

	ErrInvalidAccessToken = response.ErrorModel{
		Message:    "Please provide a valid access token",
		StatusCode: http.StatusUnauthorized,
	}

	ErrInvalidRefreshToken = response.ErrorModel{
		Message:    "Please provide a valid refresh token",
		StatusCode: http.StatusUnauthorized,
	}

	ErrSessionExpired = response.ErrorModel{
		Message:    "Your session has expired, please login again",
		StatusCode: http.StatusUnauthorized,
	}
)
