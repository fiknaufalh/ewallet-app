package middleware

import (
	"errors"
	"net/http"

	"ewallet-app/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle the first error
		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			var statusCode int
			var response gin.H

			switch {
			case errors.Is(err, entity.ErrInsufficientBalance):
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "insufficient balance"}
			
			case errors.Is(err, entity.ErrInvalidAmount):
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "invalid amount"}
			
			default:
				statusCode = http.StatusInternalServerError
				response = gin.H{"error": "internal server error"}
			}

			c.JSON(statusCode, response)
			return
		}
	}
}