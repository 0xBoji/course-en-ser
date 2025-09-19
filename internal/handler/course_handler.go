package handler

import (
	"net/http"
	"net/url"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CourseHandler handles course-related HTTP requests
type CourseHandler struct {
	courseService service.CourseService
	s3Service     *service.S3Service
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(courseService service.CourseService, s3Service *service.S3Service) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		s3Service:     s3Service,
	}
}

// CreateCourseWithImage creates a new course with image upload
// @Summary Create a new course with image upload
// @Description Create a new course with title, description, difficulty level, and optional image file
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Course title"
// @Param description formData string true "Course description"
// @Param difficulty formData string true "Course difficulty (Beginner, Intermediate, Advanced)"
// @Param image formData file false "Course image file (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 201 {object} models.CourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses [post]
func (h *CourseHandler) CreateCourseWithImage(c *gin.Context) {
	// Get form data
	title := c.PostForm("title")
	description := c.PostForm("description")
	difficulty := c.PostForm("difficulty")

	// Validate required fields
	if title == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Title is required",
		})
		return
	}

	if description == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Description is required",
		})
		return
	}

	// Validate difficulty
	validDifficulties := map[string]bool{
		"Beginner":     true,
		"Intermediate": true,
		"Advanced":     true,
	}
	if !validDifficulties[difficulty] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "Difficulty must be one of: Beginner, Intermediate, Advanced",
		})
		return
	}

	// Handle image upload (optional)
	var imageURL *string
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		// Upload image to S3
		uploadedURL, uploadErr := h.s3Service.UploadCourseImage(file)
		if uploadErr != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Image upload failed",
				Message: uploadErr.Error(),
			})
			return
		}
		imageURL = &uploadedURL
	}

	// Create course request
	req := models.CourseRequest{
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
		ImageURL:    imageURL,
	}

	course, err := h.courseService.CreateCourse(req)
	if err != nil {
		// If course creation fails and we uploaded an image, clean it up
		if imageURL != nil {
			h.s3Service.DeleteCourseImage(*imageURL)
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create course",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// CreateCourse creates a new course (JSON endpoint for backward compatibility)
// @Summary Create a new course (JSON)
// @Description Create a new course with title, description, and difficulty level using JSON
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.CourseRequest true "Course data"
// @Success 201 {object} models.CourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses/json [post]
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

	// Validate image URL if provided
	if req.ImageURL != nil && *req.ImageURL != "" {
		if !isValidURL(*req.ImageURL) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: "Image URL must be a valid URL",
			})
			return
		}
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

// isValidURL checks if a string is a valid URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
