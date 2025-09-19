package service

import (
	"errors"
	"testing"

	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEnrollmentRepository is a mock implementation of EnrollmentRepository
type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) Create(enrollment *models.Enrollment) error {
	args := m.Called(enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) GetByStudentEmail(email string) ([]models.Enrollment, error) {
	args := m.Called(email)
	return args.Get(0).([]models.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetByStudentAndCourse(email string, courseID uuid.UUID) (*models.Enrollment, error) {
	args := m.Called(email, courseID)
	return args.Get(0).(*models.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) ExistsByStudentAndCourse(email string, courseID uuid.UUID) (bool, error) {
	args := m.Called(email, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEnrollmentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestEnrollmentService_EnrollStudent(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	courseID := uuid.New()
	course := &models.Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	req := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     courseID,
	}

	enrollment := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: req.StudentEmail,
		CourseID:     req.CourseID,
	}

	mockCourseRepo.On("GetByID", courseID).Return(course, nil)
	mockEnrollmentRepo.On("Create", mock.AnythingOfType("*models.Enrollment")).Return(nil)
	mockEnrollmentRepo.On("GetByStudentAndCourse", req.StudentEmail, req.CourseID).Return(enrollment, nil)

	result, err := service.EnrollStudent(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.StudentEmail, result.StudentEmail)
	assert.Equal(t, req.CourseID, result.CourseID)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}

func TestEnrollmentService_EnrollStudent_InvalidEmail(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	req := models.EnrollmentRequest{
		StudentEmail: "invalid-email",
		CourseID:     uuid.New(),
	}

	result, err := service.EnrollStudent(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid email format", err.Error())
}

func TestEnrollmentService_EnrollStudent_CourseNotFound(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	courseID := uuid.New()
	req := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     courseID,
	}

	mockCourseRepo.On("GetByID", courseID).Return((*models.Course)(nil), gorm.ErrRecordNotFound)

	result, err := service.EnrollStudent(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "course not found", err.Error())
	mockCourseRepo.AssertExpectations(t)
}

func TestEnrollmentService_EnrollStudent_DatabaseError(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	courseID := uuid.New()
	course := &models.Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	req := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     courseID,
	}

	mockCourseRepo.On("GetByID", courseID).Return(course, nil)
	mockEnrollmentRepo.On("Create", mock.AnythingOfType("*models.Enrollment")).Return(errors.New("database error"))

	result, err := service.EnrollStudent(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}

func TestEnrollmentService_GetStudentEnrollments(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	studentEmail := "student@example.com"
	courseID1 := uuid.New()
	courseID2 := uuid.New()

	enrollments := []models.Enrollment{
		{
			ID:           uuid.New(),
			StudentEmail: studentEmail,
			CourseID:     courseID1,
			Course: models.Course{
				ID:          courseID1,
				Title:       "Course 1",
				Description: "Description 1",
				Difficulty:  "Beginner",
			},
		},
		{
			ID:           uuid.New(),
			StudentEmail: studentEmail,
			CourseID:     courseID2,
			Course: models.Course{
				ID:          courseID2,
				Title:       "Course 2",
				Description: "Description 2",
				Difficulty:  "Intermediate",
			},
		},
	}

	mockEnrollmentRepo.On("GetByStudentEmail", studentEmail).Return(enrollments, nil)

	result, err := service.GetStudentEnrollments(studentEmail)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, studentEmail, result.StudentEmail)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Enrollments, 2)
	assert.Equal(t, "Course 1", result.Enrollments[0].Course.Title)
	assert.Equal(t, "Course 2", result.Enrollments[1].Course.Title)
	mockEnrollmentRepo.AssertExpectations(t)
}

func TestEnrollmentService_GetStudentEnrollments_InvalidEmail(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	result, err := service.GetStudentEnrollments("invalid-email")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid email format", err.Error())
}

func TestEnrollmentService_GetStudentEnrollments_DatabaseError(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	studentEmail := "student@example.com"

	mockEnrollmentRepo.On("GetByStudentEmail", studentEmail).Return([]models.Enrollment(nil), errors.New("database error"))

	result, err := service.GetStudentEnrollments(studentEmail)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockEnrollmentRepo.AssertExpectations(t)
}

func TestEnrollmentService_GetStudentEnrollments_Empty(t *testing.T) {
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	service := NewEnrollmentService(mockEnrollmentRepo, mockCourseRepo)

	studentEmail := "student@example.com"
	enrollments := []models.Enrollment{}

	mockEnrollmentRepo.On("GetByStudentEmail", studentEmail).Return(enrollments, nil)

	result, err := service.GetStudentEnrollments(studentEmail)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, studentEmail, result.StudentEmail)
	assert.Equal(t, 0, result.Total)
	assert.Len(t, result.Enrollments, 0)
	mockEnrollmentRepo.AssertExpectations(t)
}
