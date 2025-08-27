package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrInternalServerError = response.ErrorModel{
		Message:    "Internal server error, Please try again later",
		StatusCode: http.StatusInternalServerError,
	}

	ErrLoginRequired = response.ErrorModel{
		Message:    "Please login to access this resource",
		StatusCode: http.StatusUnauthorized,
	}

	ErrBadRequest = response.ErrorModel{
		Message:    "Bad request, Please check your request",
		StatusCode: http.StatusBadRequest,
	}

	ErrAtleastOneField = response.ErrorModel{
		Message:    "At least one field must be provided",
		StatusCode: http.StatusBadRequest,
	}
)
