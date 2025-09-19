package repository

import (
	"errors"
	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentRepository defines the interface for enrollment data operations
type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	GetByStudentEmail(email string) ([]models.Enrollment, error)
	GetByStudentAndCourse(email string, courseID uuid.UUID) (*models.Enrollment, error)
	ExistsByStudentAndCourse(email string, courseID uuid.UUID) (bool, error)
	Delete(id uuid.UUID) error
}

// enrollmentRepository implements EnrollmentRepository interface
type enrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	exists, err := r.ExistsByStudentAndCourse(enrollment.StudentEmail, enrollment.CourseID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("student is already enrolled in this course")
	}

	return r.db.Create(enrollment).Error
}

// GetByStudentEmail retrieves all enrollments for a student
func (r *enrollmentRepository) GetByStudentEmail(email string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Preload("Course").Where("student_email = ?", email).Order("enrolled_at DESC").Find(&enrollments).Error
	return enrollments, err
}

// GetByStudentAndCourse retrieves a specific enrollment
func (r *enrollmentRepository) GetByStudentAndCourse(email string, courseID uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("Course").Where("student_email = ? AND course_id = ?", email, courseID).First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// ExistsByStudentAndCourse checks if an enrollment exists
func (r *enrollmentRepository) ExistsByStudentAndCourse(email string, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).Where("student_email = ? AND course_id = ?", email, courseID).Count(&count).Error
	return count > 0, err
}

// Delete deletes an enrollment by ID
func (r *enrollmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Enrollment{}, id).Error
}
