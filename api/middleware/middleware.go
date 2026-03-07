package middleware

import (
	"cbs-simulator/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware logs all requests using structured logger
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log after request
		latency := time.Since(startTime)
		
		// Use structured logger if available
		if utils.AppLogger != nil {
			utils.AppLogger.LogAPI(
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				c.Writer.Status(),
				latency,
			)
		}
	}
}

// AuthMiddleware validates authentication (placeholder for JWT)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In production, validate JWT token here
		// For simulator, we skip authentication
		c.Next()
	}
}
