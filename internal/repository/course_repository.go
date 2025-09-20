package repository

import (
	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CourseRepository defines the interface for course data operations
type CourseRepository interface {
	Create(course *models.Course) error
	GetAll() ([]models.Course, error)
	GetWithPagination(params models.CourseQueryParams) ([]models.Course, int, error)
	GetByID(id uuid.UUID) (*models.Course, error)
	Update(course *models.Course) error
	Delete(id uuid.UUID) error
	ExistsByID(id uuid.UUID) (bool, error)
}

// courseRepository implements CourseRepository interface
type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

// Create creates a new course
func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// GetAll retrieves all courses (backward compatibility)
func (r *courseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Order("created_at DESC").Find(&courses).Error
	return courses, err
}

// GetWithPagination retrieves courses with pagination, search, and filtering
func (r *courseRepository) GetWithPagination(params models.CourseQueryParams) ([]models.Course, int, error) {
	var courses []models.Course
	var totalCount int64

	// Build base query
	query := r.db.Model(&models.Course{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	// Apply difficulty filter
	if len(params.Difficulty) > 0 {
		query = query.Where("difficulty IN ?", params.Difficulty)
	}

	// Get total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int(totalCount), nil
}

// GetByID retrieves a course by ID
func (r *courseRepository) GetByID(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	err := r.db.Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// Update updates an existing course
func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete deletes a course by ID
func (r *courseRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Course{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ExistsByID checks if a course exists by ID
func (r *courseRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Course{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
