package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorModel struct {
	Message    string
	StatusCode int
}

func (e ErrorModel) Error() string {
	return e.Message
}

func (e ErrorModel) GetStatusCode() int {
	return e.StatusCode
}

func (e ErrorModel) GetMessage() string {
	return e.Message
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

func BuildErrorResponse(c *gin.Context, err interface{}) {
	var message string
	var statusCode int

	if errorModel, ok := err.(ErrorModel); ok {
		message = errorModel.GetMessage()
		statusCode = errorModel.GetStatusCode()
	} else if errorErr, ok := err.(error); ok {
		message = errorErr.Error()
		statusCode = http.StatusInternalServerError
	} else {
		message = "Something went wrong. Please try again later"
		statusCode = http.StatusInternalServerError
	}

	response := ErrorResponse{
		Success: false,
		Message: message,
	}

	c.JSON(statusCode, response)
}

func BuildValidationErrorResponse(c *gin.Context, validationErrors map[string]string) {
	message := "Please check your input and try again"
	if len(validationErrors) > 1 {
		message = "Some required information is missing or incorrect"
	}

	response := ErrorResponse{
		Success: false,
		Message: message,
		Error:   validationErrors,
	}

	c.JSON(http.StatusBadRequest, response)
}
