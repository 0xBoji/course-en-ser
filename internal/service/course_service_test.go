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

// MockCourseRepository is a mock implementation of CourseRepository
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) Create(course *models.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) GetAll() ([]models.Course, error) {
	args := m.Called()
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCourseRepository) GetWithPagination(params models.CourseQueryParams) ([]models.Course, int, error) {
	args := m.Called(params)
	return args.Get(0).([]models.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) GetByID(id uuid.UUID) (*models.Course, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseRepository) Update(course *models.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCourseService_CreateCourse(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	req := models.CourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Course")).Return(nil)

	result, err := service.CreateCourse(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Description, result.Description)
	assert.Equal(t, req.Difficulty, result.Difficulty)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_CreateCourse_Error(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	req := models.CourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Course")).Return(errors.New("database error"))

	result, err := service.CreateCourse(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetAllCourses(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	courses := []models.Course{
		{
			ID:          uuid.New(),
			Title:       "Course 1",
			Description: "Description 1",
			Difficulty:  "Beginner",
		},
		{
			ID:          uuid.New(),
			Title:       "Course 2",
			Description: "Description 2",
			Difficulty:  "Advanced",
		},
	}

	mockRepo.On("GetAll").Return(courses, nil)

	result, err := service.GetAllCourses()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, courses[0].Title, result[0].Title)
	assert.Equal(t, courses[1].Title, result[1].Title)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCoursesWithPagination(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	courses := []models.Course{
		{
			ID:          uuid.New(),
			Title:       "Go Programming",
			Description: "Learn Go programming",
			Difficulty:  "Beginner",
		},
		{
			ID:          uuid.New(),
			Title:       "Advanced Go",
			Description: "Advanced Go concepts",
			Difficulty:  "Advanced",
		},
	}

	params := models.CourseQueryParams{
		Page:       1,
		Limit:      10,
		Search:     "Go",
		Difficulty: []string{"Beginner", "Advanced"},
	}

	mockRepo.On("GetWithPagination", params).Return(courses, 2, nil)

	result, err := service.GetCoursesWithPagination(params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Pagination.TotalCount)
	assert.Equal(t, 1, result.Pagination.CurrentPage)
	assert.Equal(t, 1, result.Pagination.TotalPages)
	assert.False(t, result.Pagination.HasNext)
	assert.False(t, result.Pagination.HasPrev)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCoursesWithPagination_DefaultValues(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil)

	courses := []models.Course{}

	// Test with invalid/empty parameters
	params := models.CourseQueryParams{
		Page:  0, // Should default to 1
		Limit: 0, // Should default to 10
	}

	expectedParams := models.CourseQueryParams{
		Page:       1,
		Limit:      10,
		Search:     "",
		Difficulty: nil,
	}

	mockRepo.On("GetWithPagination", expectedParams).Return(courses, 0, nil)

	result, err := service.GetCoursesWithPagination(params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 0)
	assert.Equal(t, 0, result.Pagination.TotalCount)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	courseID := uuid.New()
	course := &models.Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	mockRepo.On("GetByID", courseID).Return(course, nil)

	result, err := service.GetCourseByID(courseID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, course.ID, result.ID)
	assert.Equal(t, course.Title, result.Title)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID_NotFound(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo, nil) // No Redis for unit tests

	courseID := uuid.New()

	mockRepo.On("GetByID", courseID).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetCourseByID(courseID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "course not found", err.Error())
	mockRepo.AssertExpectations(t)
}
