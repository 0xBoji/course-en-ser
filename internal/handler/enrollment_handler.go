package handler

import (
	"net/http"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EnrollmentHandler handles enrollment-related HTTP requests
type EnrollmentHandler struct {
	enrollmentService service.EnrollmentService
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(enrollmentService service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

// EnrollStudent enrolls a student in a course
// @Summary Enroll a student in a course
// @Description Enroll a student in a specific course using their email and course ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param enrollment body models.EnrollmentRequest true "Enrollment data"
// @Success 201 {object} models.EnrollmentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /enrollments [post]
func (h *EnrollmentHandler) EnrollStudent(c *gin.Context) {
	var req models.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if req.StudentEmail == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Student email is required",
		})
		return
	}

	if req.CourseID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Course ID is required",
		})
		return
	}

	enrollment, err := h.enrollmentService.EnrollStudent(req)
	if err != nil {
		if err.Error() == "invalid email format" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: "Invalid email format",
			})
			return
		}
		if err.Error() == "course not found" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Course not found",
				Message: "The specified course does not exist",
			})
			return
		}
		if err.Error() == "student is already enrolled in this course" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Enrollment conflict",
				Message: "Student is already enrolled in this course",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to enroll student",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

// GetStudentEnrollments retrieves all enrollments for a student
// @Summary Get student enrollments
// @Description Retrieve all courses a specific student is enrolled in
// @Tags enrollments
// @Produce json
// @Param email path string true "Student email"
// @Success 200 {object} models.StudentEnrollmentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /students/{email}/enrollments [get]
func (h *EnrollmentHandler) GetStudentEnrollments(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Student email is required",
		})
		return
	}

	enrollments, err := h.enrollmentService.GetStudentEnrollments(email)
	if err != nil {
		if err.Error() == "invalid email format" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: "Invalid email format",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve enrollments",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}
