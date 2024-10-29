package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const IdempotencyKeyHeader = "X-Idempotency-Key"

func RequireIdempotencyKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip idempotency check for GET requests
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader(IdempotencyKeyHeader)
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Idempotency-Key header is required for this request",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}