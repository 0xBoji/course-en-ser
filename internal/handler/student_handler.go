package handler

import (
	"log"
	"net/http"
	"sonic-labs/course-enrollment-service/internal/constants"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentHandler handles student-related HTTP requests
type StudentHandler struct {
	studentService service.StudentService
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(studentService service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// GetAllStudents retrieves all students with enrollment statistics
// @Summary Get all students
// @Description Get all students with their enrollment count and statistics (Admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} models.AllStudentsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/students [get]
func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	log.Printf("API Request: GET %s from %s", c.Request.URL.Path, c.ClientIP())

	response, err := h.studentService.GetAllStudents()
	if err != nil {
		log.Printf("API Response: GET %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to retrieve students",
		})
		return
	}

	log.Printf("API Response: GET %s -> 200", c.Request.URL.Path)
	c.JSON(http.StatusOK, response)
}

// GetAllEnrollments retrieves all enrollments with course details
// @Summary Get all enrollments
// @Description Get all enrollments with course details (Admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} models.AllEnrollmentsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/enrollments [get]
func (h *StudentHandler) GetAllEnrollments(c *gin.Context) {
	log.Printf("API Request: GET %s from %s", c.Request.URL.Path, c.ClientIP())

	response, err := h.studentService.GetAllEnrollments()
	if err != nil {
		log.Printf("API Response: GET %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to retrieve enrollments",
		})
		return
	}

	log.Printf("API Response: GET %s -> 200", c.Request.URL.Path)
	c.JSON(http.StatusOK, response)
}

// DeleteEnrollment deletes an enrollment
// @Summary Delete an enrollment
// @Description Delete an enrollment by ID (Admin only)
// @Tags admin
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/enrollments/{id} [delete]
func (h *StudentHandler) DeleteEnrollment(c *gin.Context) {
	log.Printf("API Request: DELETE %s from %s", c.Request.URL.Path, c.ClientIP())

	// Parse enrollment ID
	enrollmentIDStr := c.Param("id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		log.Printf("API Response: DELETE %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid enrollment ID format",
		})
		return
	}

	// Delete enrollment
	err = h.studentService.DeleteEnrollment(enrollmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("API Response: DELETE %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Enrollment not found",
			})
			return
		}

		log.Printf("API Response: DELETE %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to delete enrollment",
		})
		return
	}

	log.Printf("API Response: DELETE %s -> 204", c.Request.URL.Path)
	c.Status(http.StatusNoContent)
}
