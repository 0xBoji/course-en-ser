package handler

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Validation failed"`
	Message string `json:"message" example:"Title is required"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}
