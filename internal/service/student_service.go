package service

import (
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/repository"

	"github.com/google/uuid"
)

// StudentService defines the interface for student business logic
type StudentService interface {
	GetAllStudents() (*models.AllStudentsResponse, error)
	GetAllEnrollments() (*models.AllEnrollmentsResponse, error)
	DeleteEnrollment(id uuid.UUID) error
}

// studentService implements StudentService interface
type studentService struct {
	enrollmentRepo repository.EnrollmentRepository
}

// NewStudentService creates a new student service
func NewStudentService(enrollmentRepo repository.EnrollmentRepository) StudentService {
	return &studentService{
		enrollmentRepo: enrollmentRepo,
	}
}

// GetAllStudents retrieves all students with their enrollment statistics
func (s *studentService) GetAllStudents() (*models.AllStudentsResponse, error) {
	students, err := s.enrollmentRepo.GetAllStudents()
	if err != nil {
		return nil, err
	}

	return &models.AllStudentsResponse{
		Students: students,
		Total:    len(students),
	}, nil
}

// GetAllEnrollments retrieves all enrollments with course details
func (s *studentService) GetAllEnrollments() (*models.AllEnrollmentsResponse, error) {
	enrollments, err := s.enrollmentRepo.GetAllEnrollments()
	if err != nil {
		return nil, err
	}

	return &models.AllEnrollmentsResponse{
		Enrollments: enrollments,
		Total:       len(enrollments),
	}, nil
}

// DeleteEnrollment deletes an enrollment by ID
func (s *studentService) DeleteEnrollment(id uuid.UUID) error {
	return s.enrollmentRepo.Delete(id)
}
