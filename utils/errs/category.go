package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrCategoryNotFound = response.ErrorModel{
		Message:    "Ticket category not found",
		StatusCode: http.StatusNotFound,
	}

	ErrCategoryAlreadyExists = response.ErrorModel{
		Message:    "Ticket category already exists",
		StatusCode: http.StatusBadRequest,
	}

	ErrStockNotAvailable = response.ErrorModel{
		Message:    "Stock not available",
		StatusCode: http.StatusBadRequest,
	}
)
