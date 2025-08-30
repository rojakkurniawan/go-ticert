package errs

import (
	"net/http"
	"ticert/utils/response"
)

var (
	ErrOrderNotFound = response.ErrorModel{
		Message:    "Order not found",
		StatusCode: http.StatusNotFound,
	}
	ErrQuantityNotMatch = response.ErrorModel{
		Message:    "The number of order details requested does not match the quantity",
		StatusCode: http.StatusBadRequest,
	}

	ErrOrderAlreadyCancelled = response.ErrorModel{
		Message:    "Order already cancelled",
		StatusCode: http.StatusBadRequest,
	}

	ErrOrderAlreadyPaid = response.ErrorModel{
		Message:    "Order already paid",
		StatusCode: http.StatusBadRequest,
	}

	ErrTicketAlreadyRedeemed = response.ErrorModel{
		Message:    "Ticket already redeemed",
		StatusCode: http.StatusBadRequest,
	}

	ErrTicketNotFound = response.ErrorModel{
		Message:    "Ticket not found",
		StatusCode: http.StatusNotFound,
	}

	ErrOrderNotPaid = response.ErrorModel{
		Message:    "Order not paid",
		StatusCode: http.StatusBadRequest,
	}
)
