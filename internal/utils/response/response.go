package response

import (
	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Paginated(c *gin.Context, message string, data interface{}, pagination models.Pagination) {
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		c.JSON(appErr.StatusCode, models.StandardResponse{
			Success: false,
			Message: appErr.Message,
			Error:   appErr.Error(),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, models.StandardResponse{
		Success: false,
		Message: "Internal server error",
		Error:   err.Error(),
	})
}

func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, models.StandardResponse{
		Success: false,
		Message: "Invalid input parameters",
		Error:   err.Error(),
	})
}
