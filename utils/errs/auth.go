package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrEmailAlreadyExists = response.ErrorModel{
		Message:    "Email already exists",
		StatusCode: http.StatusBadRequest,
	}

	ErrAuthInvalidCredentials = response.ErrorModel{
		Message:    "Email or password is incorrect, Please try again",
		StatusCode: http.StatusUnauthorized,
	}

	ErrUserNotFound = response.ErrorModel{
		Message:    "User not found",
		StatusCode: http.StatusNotFound,
	}
)
