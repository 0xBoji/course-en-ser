package service

import (
	"errors"
	"net/mail"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentService defines the interface for enrollment business logic
type EnrollmentService interface {
	EnrollStudent(req models.EnrollmentRequest) (*models.EnrollmentResponse, error)
	GetStudentEnrollments(email string) (*models.StudentEnrollmentsResponse, error)
	UnenrollStudent(email string, courseID uuid.UUID) error
}

// enrollmentService implements EnrollmentService interface
type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, courseRepo repository.CourseRepository) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
	}
}

func (s *enrollmentService) EnrollStudent(req models.EnrollmentRequest) (*models.EnrollmentResponse, error) {
	if _, err := mail.ParseAddress(req.StudentEmail); err != nil {
		return nil, errors.New("invalid email format")
	}
	_, err := s.courseRepo.GetByID(req.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	enrollment := models.Enrollment{
		StudentEmail: req.StudentEmail,
		CourseID:     req.CourseID,
	}

	if err := s.enrollmentRepo.Create(&enrollment); err != nil {
		return nil, err
	}
	createdEnrollment, err := s.enrollmentRepo.GetByStudentAndCourse(req.StudentEmail, req.CourseID)
	if err != nil {
		return nil, err
	}

	response := createdEnrollment.ToResponse()
	return &response, nil
}

func (s *enrollmentService) GetStudentEnrollments(email string) (*models.StudentEnrollmentsResponse, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}

	enrollments, err := s.enrollmentRepo.GetByStudentEmail(email)
	if err != nil {
		return nil, err
	}

	responses := make([]models.EnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = enrollment.ToResponse()
	}

	return &models.StudentEnrollmentsResponse{
		StudentEmail: email,
		Enrollments:  responses,
		Total:        len(responses),
	}, nil
}

func (s *enrollmentService) UnenrollStudent(email string, courseID uuid.UUID) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}

	enrollment, err := s.enrollmentRepo.GetByStudentAndCourse(email, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	return s.enrollmentRepo.Delete(enrollment.ID)
}
