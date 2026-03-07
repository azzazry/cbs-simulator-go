package middleware

import (
	"bytes"
	"cbs-simulator/models"
	"cbs-simulator/services"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
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

// LoggerMiddleware logs all requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		gin.DefaultWriter.Write([]byte(
			time.Now().Format("2006/01/02 - 15:04:05") +
				" | " + http.StatusText(statusCode) +
				" | " + latency.String() +
				" | " + clientIP +
				" | " + method + " " + path + "\n",
		))
	}
}

// AuthMiddleware validates JWT tokens on protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header format must be: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Set claims in context for downstream handlers
		c.Set("cif", claims.CIF)
		c.Set("role", claims.Role)
		c.Set("token_type", claims.TokenType)
		c.Set("jti", claims.ID)
		if claims.ExpiresAt != nil {
			c.Set("token_expires_at", claims.ExpiresAt.Time)
		}

		c.Next()
	}
}

// RequireRole checks if the authenticated user has the required role
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cif, exists := c.Get("cif")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		for _, requiredRole := range requiredRoles {
			hasRole, err := services.HasRole(cif.(string), requiredRole)
			if err == nil && hasRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Insufficient permissions. Required role: " + strings.Join(requiredRoles, " or "),
		})
		c.Abort()
	}
}

// AuditMiddleware logs all API activity for security compliance
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read request body (need to save and restore for handlers)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Sanitize body (remove sensitive fields)
		bodyStr := sanitizeRequestBody(string(bodyBytes))

		c.Next()

		// Get CIF from context (set by AuthMiddleware)
		cif := ""
		if cifVal, exists := c.Get("cif"); exists {
			cif = cifVal.(string)
		}

		// Log audit entry
		auditLog := models.AuditLog{
			CIF:            cif,
			Action:         c.Request.Method + " " + c.Request.URL.Path,
			Resource:       c.Request.URL.Path,
			IPAddress:      c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
			RequestMethod:  c.Request.Method,
			RequestPath:    c.Request.URL.Path,
			RequestBody:    bodyStr,
			ResponseStatus: c.Writer.Status(),
		}

		// Log asynchronously to not block the response
		go services.LogAudit(auditLog)
	}
}

// sanitizeRequestBody removes sensitive fields from request body for audit logging
func sanitizeRequestBody(body string) string {
	if len(body) > 1000 {
		return body[:1000] + "...[truncated]"
	}
	// Simple sanitization: mask PIN and password fields
	body = strings.ReplaceAll(body, "\"pin\":", "\"pin\":\"***\",\"_pin\":")
	return body
}

// RateLimiterMiddleware implements a simple in-memory rate limiter
func RateLimiterMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	type clientInfo struct {
		count   int
		resetAt time.Time
	}

	var mu sync.Mutex
	clients := make(map[string]*clientInfo)

	// Cleanup old entries periodically
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			now := time.Now()
			for key, info := range clients {
				if now.After(info.resetAt) {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.Lock()
		info, exists := clients[clientIP]
		now := time.Now()

		if !exists || now.After(info.resetAt) {
			clients[clientIP] = &clientInfo{
				count:   1,
				resetAt: now.Add(window),
			}
			mu.Unlock()
			c.Next()
			return
		}

		info.count++
		if info.count > maxRequests {
			mu.Unlock()
			remaining := info.resetAt.Sub(now).Seconds()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":      "error",
				"message":     "Rate limit exceeded. Too many requests.",
				"retry_after": int(remaining),
			})
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}
