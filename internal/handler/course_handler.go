package handler

import (
	"net/http"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CourseHandler handles course-related HTTP requests
type CourseHandler struct {
	courseService service.CourseService
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// CreateCourse creates a new course
// @Summary Create a new course
// @Description Create a new course with title, description, and difficulty level
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.CourseRequest true "Course data"
// @Success 201 {object} models.CourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Title is required",
		})
		return
	}

	if req.Description == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Description is required",
		})
		return
	}
	validDifficulties := map[string]bool{
		"Beginner":     true,
		"Intermediate": true,
		"Advanced":     true,
	}
	if !validDifficulties[req.Difficulty] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Difficulty must be one of: Beginner, Intermediate, Advanced",
		})
		return
	}

	course, err := h.courseService.CreateCourse(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create course",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetAllCourses retrieves all courses
// @Summary Get all courses
// @Description Retrieve a list of all available courses
// @Tags courses
// @Produce json
// @Success 200 {array} models.CourseResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve courses",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// GetCourseByID retrieves a course by ID
// @Summary Get course by ID
// @Description Retrieve a specific course by its ID
// @Tags courses
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid course ID",
			Message: "Course ID must be a valid UUID",
		})
		return
	}

	course, err := h.courseService.GetCourseByID(id)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Course not found",
				Message: "The requested course does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve course",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, course)
}
