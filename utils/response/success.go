package response

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func BuildSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination *Pagination) {
	var response SuccessResponse
	response.Success = true
	response.Message = message
	response.Data = data
	response.Pagination = pagination
	c.JSON(statusCode, response)
}
