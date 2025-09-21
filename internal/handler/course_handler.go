package handler

import (
	"log"
	"net/http"
	"net/url"
	"sonic-labs/course-enrollment-service/internal/constants"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

// GetAllCourses retrieves all courses with pagination and search
// @Summary Get all courses with pagination and search
// @Description Retrieve a list of courses with optional pagination, search, and filtering
// @Tags courses
// @Produce json
// @Param page query int false "Page number (default: 1)" example(1)
// @Param limit query int false "Items per page (default: 10, max: 100)" example(10)
// @Param search query string false "Search in title and description" example("golang")
// @Param difficulty query []string false "Filter by difficulty levels" example("Beginner,Intermediate")
// @Success 200 {object} models.CourseListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	// Parse query parameters
	var params models.CourseQueryParams

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	// Parse search
	params.Search = strings.TrimSpace(c.Query("search"))

	// Parse difficulty filter
	if difficultyStr := c.Query("difficulty"); difficultyStr != "" {
		difficulties := strings.Split(difficultyStr, ",")
		validDifficulties := []string{}
		for _, d := range difficulties {
			d = strings.TrimSpace(d)
			if d == "Beginner" || d == "Intermediate" || d == "Advanced" {
				validDifficulties = append(validDifficulties, d)
			}
		}
		params.Difficulty = validDifficulties
	}

	// Check if any pagination/search parameters are provided
	hasPaginationParams := params.Page > 0 || params.Limit > 0 || params.Search != "" || len(params.Difficulty) > 0

	if hasPaginationParams {
		// Use new pagination endpoint
		result, err := h.courseService.GetCoursesWithPagination(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to retrieve courses",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
	} else {
		// Backward compatibility: return simple array for existing clients
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

// UpdateCourse updates an existing course
// @Summary Update a course
// @Description Update an existing course by ID (Admin only)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param course body models.CourseRequest true "Course update data"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	log.Printf("API Request: PUT %s from %s", c.Request.URL.Path, c.ClientIP())

	// Parse course ID
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		log.Printf("API Response: PUT %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid course ID format",
		})
		return
	}

	// Parse request body
	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("API Response: PUT %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update course
	response, err := h.courseService.UpdateCourse(courseID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("API Response: PUT %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Course not found",
			})
			return
		}

		log.Printf("API Response: PUT %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to update course",
		})
		return
	}

	log.Printf("API Response: PUT %s -> 200", c.Request.URL.Path)
	c.JSON(http.StatusOK, response)
}

// DeleteCourse deletes a course
// @Summary Delete a course
// @Description Delete a course by ID (Admin only)
// @Tags courses
// @Produce json
// @Param id path string true "Course ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	log.Printf("API Request: DELETE %s from %s", c.Request.URL.Path, c.ClientIP())

	// Parse course ID
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		log.Printf("API Response: DELETE %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid course ID format",
		})
		return
	}

	// Delete course
	err = h.courseService.DeleteCourse(courseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("API Response: DELETE %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Course not found",
			})
			return
		}

		log.Printf("API Response: DELETE %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to delete course",
		})
		return
	}

	log.Printf("API Response: DELETE %s -> 204", c.Request.URL.Path)
	c.Status(http.StatusNoContent)
}

// GetCourseStudents retrieves all students enrolled in a specific course
// @Summary Get course students
// @Description Get all student emails enrolled in a specific course (Admin only)
// @Tags admin
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]interface{} "{"students": ["email1", "email2"], "total": 2}"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses/{id}/students [get]
func (h *CourseHandler) GetCourseStudents(c *gin.Context) {
	log.Printf("API Request: GET %s from %s", c.Request.URL.Path, c.ClientIP())

	// Parse course ID
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		log.Printf("API Response: GET %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid course ID format",
		})
		return
	}

	// Get course students
	students, err := h.courseService.GetCourseStudents(courseID)
	if err != nil {
		if err.Error() == "course not found" {
			log.Printf("API Response: GET %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Course not found",
			})
			return
		}

		log.Printf("API Response: GET %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to retrieve course students",
		})
		return
	}

	log.Printf("API Response: GET %s -> 200", c.Request.URL.Path)
	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

// RemoveStudentFromCourse removes a student from a specific course
// @Summary Remove student from course
// @Description Remove a student from a specific course (Admin only)
// @Tags admin
// @Produce json
// @Param id path string true "Course ID"
// @Param email path string true "Student Email"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /courses/{id}/students/{email} [delete]
func (h *CourseHandler) RemoveStudentFromCourse(c *gin.Context) {
	log.Printf("API Request: DELETE %s from %s", c.Request.URL.Path, c.ClientIP())

	// Parse course ID
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		log.Printf("API Response: DELETE %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Invalid course ID format",
		})
		return
	}

	// Get student email
	studentEmail := c.Param("email")
	if studentEmail == "" {
		log.Printf("API Response: DELETE %s -> 400", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   constants.HTTPBadRequest,
			Message: "Student email is required",
		})
		return
	}

	// Remove student from course
	err = h.courseService.RemoveStudentFromCourse(courseID, studentEmail)
	if err != nil {
		if err.Error() == "course not found" {
			log.Printf("API Response: DELETE %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Course not found",
			})
			return
		}
		if err.Error() == "student not enrolled in this course" {
			log.Printf("API Response: DELETE %s -> 404", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   constants.HTTPNotFound,
				Message: "Student not enrolled in this course",
			})
			return
		}

		log.Printf("API Response: DELETE %s -> 500", c.Request.URL.Path)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   constants.HTTPInternalServerError,
			Message: "Failed to remove student from course",
		})
		return
	}

	log.Printf("API Response: DELETE %s -> 204", c.Request.URL.Path)
	c.Status(http.StatusNoContent)
}

// isValidURL checks if a string is a valid URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
