package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrEventNotFound = response.ErrorModel{
		Message:    "Event not found",
		StatusCode: http.StatusNotFound,
	}

	ErrEventTitleAlreadyExists = response.ErrorModel{
		Message:    "Event with this title already exists",
		StatusCode: http.StatusBadRequest,
	}
)
