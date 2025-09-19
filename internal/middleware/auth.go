package middleware

import (
	"net/http"
	"strings"

	"sonic-labs/course-enrollment-service/internal/auth"
	"sonic-labs/course-enrollment-service/internal/constants"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// AuthMiddleware validates JWT tokens and protects routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authorization required",
				Message: constants.MsgAuthHeaderRequired,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid authorization format",
				Message: constants.MsgInvalidTokenFormat,
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Token required",
				Message: "JWT token is required",
			})
			c.Abort()
			return
		}

		// Validate the token using auth package
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid token",
				Message: constants.MsgJWTTokenInvalid,
			})
			c.Abort()
			return
		}

		// Store user information in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminMiddleware ensures the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication required",
				Message: "User role not found in context",
			})
			c.Abort()
			return
		}

		if role != constants.RoleAdmin {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Insufficient permissions",
				Message: constants.MsgAdminAccessRequired,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuthMiddleware combines authentication and admin authorization
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, validate JWT token
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authorization required",
				Message: constants.MsgAuthHeaderRequired,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid authorization format",
				Message: constants.MsgInvalidTokenFormat,
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Token required",
				Message: "JWT token is required",
			})
			c.Abort()
			return
		}

		// Validate the token using auth package
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid token",
				Message: constants.MsgJWTTokenInvalid,
			})
			c.Abort()
			return
		}

		// Check if user is admin
		if claims.Role != constants.RoleAdmin {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Insufficient permissions",
				Message: constants.MsgAdminAccessRequired,
			})
			c.Abort()
			return
		}

		// Store user information in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
