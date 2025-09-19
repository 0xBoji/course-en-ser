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
	GetByID(id uuid.UUID) (*models.Course, error)
	Update(course *models.Course) error
	Delete(id uuid.UUID) error
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

// GetAll retrieves all courses
func (r *courseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Order("created_at DESC").Find(&courses).Error
	return courses, err
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
	return r.db.Delete(&models.Course{}, id).Error
}
