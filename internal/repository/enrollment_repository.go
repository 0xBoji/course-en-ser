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
	GetAllStudents() ([]models.StudentResponse, error)
	GetAllEnrollments() ([]models.EnrollmentWithCourse, error)
	GetByID(id uuid.UUID) (*models.Enrollment, error)
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
	result := r.db.Delete(&models.Enrollment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetByID retrieves an enrollment by ID
func (r *enrollmentRepository) GetByID(id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("Course").Where("id = ?", id).First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// GetAllStudents retrieves all unique students with their enrollment count
func (r *enrollmentRepository) GetAllStudents() ([]models.StudentResponse, error) {
	var students []models.StudentResponse

	query := `
		SELECT
			student_email as email,
			COUNT(*) as enrollment_count,
			MAX(enrolled_at) as last_enrolled_at
		FROM enrollments
		GROUP BY student_email
		ORDER BY enrollment_count DESC, last_enrolled_at DESC
	`

	err := r.db.Raw(query).Scan(&students).Error
	return students, err
}

// GetAllEnrollments retrieves all enrollments with course details
func (r *enrollmentRepository) GetAllEnrollments() ([]models.EnrollmentWithCourse, error) {
	var enrollments []models.Enrollment
	err := r.db.Preload("Course").Order("enrolled_at DESC").Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	var result []models.EnrollmentWithCourse
	for _, enrollment := range enrollments {
		result = append(result, models.EnrollmentWithCourse{
			ID:           enrollment.ID,
			StudentEmail: enrollment.StudentEmail,
			Course:       enrollment.Course.ToResponse(),
			EnrolledAt:   enrollment.EnrolledAt,
		})
	}

	return result, nil
}
