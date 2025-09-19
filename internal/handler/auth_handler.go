package handler

import (
	"net/http"

	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate a user with username and password and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Username is required",
		})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Password is required",
		})
		return
	}

	// Authenticate user
	loginResponse, err := h.authService.Login(req)
	if err != nil {
		if err.Error() == "invalid username or password" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication failed",
				Message: "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user information from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Authentication required",
			Message: "User ID not found in context",
		})
		return
	}

	username, _ := c.Get("username")
	role, _ := c.Get("role")

	// Parse user ID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "Failed to parse user ID",
		})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, models.UserResponse{
		ID:       userUUID,
		Username: username.(string),
		Role:     role.(string),
	})
}
