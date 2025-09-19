package service

import (
	"errors"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CourseService defines the interface for course business logic
type CourseService interface {
	CreateCourse(req models.CourseRequest) (*models.CourseResponse, error)
	GetAllCourses() ([]models.CourseResponse, error)
	GetCourseByID(id uuid.UUID) (*models.CourseResponse, error)
	UpdateCourse(id uuid.UUID, req models.CourseRequest) (*models.CourseResponse, error)
	DeleteCourse(id uuid.UUID) error
}

// courseService implements CourseService interface
type courseService struct {
	courseRepo repository.CourseRepository
}

// NewCourseService creates a new course service
func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{
		courseRepo: courseRepo,
	}
}

func (s *courseService) CreateCourse(req models.CourseRequest) (*models.CourseResponse, error) {
	course := models.Course{
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		ImageURL:    req.ImageURL,
	}

	if err := s.courseRepo.Create(&course); err != nil {
		return nil, err
	}

	response := course.ToResponse()
	return &response, nil
}

// GetAllCourses retrieves all courses
func (s *courseService) GetAllCourses() ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToResponse()
	}

	return responses, nil
}

// GetCourseByID retrieves a course by ID
func (s *courseService) GetCourseByID(id uuid.UUID) (*models.CourseResponse, error) {
	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	response := course.ToResponse()
	return &response, nil
}

// UpdateCourse updates an existing course
func (s *courseService) UpdateCourse(id uuid.UUID, req models.CourseRequest) (*models.CourseResponse, error) {
	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	course.Title = req.Title
	course.Description = req.Description
	course.Difficulty = req.Difficulty
	course.ImageURL = req.ImageURL

	if err := s.courseRepo.Update(course); err != nil {
		return nil, err
	}

	response := course.ToResponse()
	return &response, nil
}

// DeleteCourse deletes a course
func (s *courseService) DeleteCourse(id uuid.UUID) error {
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	return s.courseRepo.Delete(id)
}
