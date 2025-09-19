package repository

import (
	"testing"

	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// EnrollmentRepositoryTestSuite defines the test suite for enrollment repository tests
type EnrollmentRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo EnrollmentRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *EnrollmentRepositoryTestSuite) SetupSuite() {
	// Initialize in-memory SQLite database for testing
	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	suite.Require().NoError(err)

	// Create tables
	err = suite.db.Exec(`
		CREATE TABLE courses (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			difficulty TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	suite.Require().NoError(err)

	err = suite.db.Exec(`
		CREATE TABLE enrollments (
			id TEXT PRIMARY KEY,
			student_email TEXT NOT NULL,
			course_id TEXT NOT NULL,
			enrolled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
			UNIQUE(student_email, course_id)
		)
	`).Error
	suite.Require().NoError(err)

	// Initialize repository
	suite.repo = NewEnrollmentRepository(suite.db)
}

// TearDownSuite runs once after all tests in the suite
func (suite *EnrollmentRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// SetupTest runs before each test
func (suite *EnrollmentRepositoryTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.db.Exec("DELETE FROM enrollments")
	suite.db.Exec("DELETE FROM courses")
}

// createTestCourse is a helper function to create a test course
func (suite *EnrollmentRepositoryTestSuite) createTestCourse(title, description, difficulty string) *models.Course {
	course := &models.Course{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
	}

	err := suite.db.Create(course).Error
	suite.Require().NoError(err)

	return course
}

// TestEnrollmentRepository_Create tests creating an enrollment
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_Create() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	enrollment := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	err := suite.repo.Create(enrollment)

	suite.NoError(err)
	suite.NotEqual(uuid.Nil, enrollment.ID)

	// Verify enrollment was created in database
	var dbEnrollment models.Enrollment
	err = suite.db.First(&dbEnrollment, "id = ?", enrollment.ID.String()).Error
	suite.NoError(err)
	suite.Equal(enrollment.StudentEmail, dbEnrollment.StudentEmail)
	suite.Equal(enrollment.CourseID, dbEnrollment.CourseID)
}

// TestEnrollmentRepository_GetByStudentEmail tests retrieving enrollments by student email
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_GetByStudentEmail() {
	course1 := suite.createTestCourse("Course 1", "Description 1", "Beginner")
	course2 := suite.createTestCourse("Course 2", "Description 2", "Intermediate")

	studentEmail := "student@example.com"

	// Create enrollments
	enrollment1 := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: studentEmail,
		CourseID:     course1.ID,
	}
	enrollment2 := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: studentEmail,
		CourseID:     course2.ID,
	}

	err := suite.repo.Create(enrollment1)
	suite.NoError(err)
	err = suite.repo.Create(enrollment2)
	suite.NoError(err)

	// Get enrollments by student email
	enrollments, err := suite.repo.GetByStudentEmail(studentEmail)

	suite.NoError(err)
	suite.Len(enrollments, 2)

	// Verify enrollment data (order might vary)
	enrollmentMap := make(map[string]models.Enrollment)
	for _, enrollment := range enrollments {
		enrollmentMap[enrollment.Course.Title] = enrollment
	}

	suite.Contains(enrollmentMap, "Course 1")
	suite.Contains(enrollmentMap, "Course 2")
	suite.Equal(studentEmail, enrollmentMap["Course 1"].StudentEmail)
	suite.Equal(studentEmail, enrollmentMap["Course 2"].StudentEmail)
}

// TestEnrollmentRepository_GetByStudentEmail_Empty tests retrieving enrollments when none exist
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_GetByStudentEmail_Empty() {
	enrollments, err := suite.repo.GetByStudentEmail("nonexistent@example.com")

	suite.NoError(err)
	suite.Len(enrollments, 0)
}

// TestEnrollmentRepository_GetByStudentAndCourse tests retrieving a specific enrollment
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_GetByStudentAndCourse() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	enrollment := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	err := suite.repo.Create(enrollment)
	suite.NoError(err)

	// Get enrollment by student and course
	retrievedEnrollment, err := suite.repo.GetByStudentAndCourse(enrollment.StudentEmail, enrollment.CourseID)

	suite.NoError(err)
	suite.NotNil(retrievedEnrollment)
	suite.Equal(enrollment.ID, retrievedEnrollment.ID)
	suite.Equal(enrollment.StudentEmail, retrievedEnrollment.StudentEmail)
	suite.Equal(enrollment.CourseID, retrievedEnrollment.CourseID)
}

// TestEnrollmentRepository_GetByStudentAndCourse_NotFound tests retrieving a non-existent enrollment
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_GetByStudentAndCourse_NotFound() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	retrievedEnrollment, err := suite.repo.GetByStudentAndCourse("nonexistent@example.com", course.ID)

	suite.Error(err)
	suite.Nil(retrievedEnrollment)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

// TestEnrollmentRepository_ExistsByStudentAndCourse tests checking if enrollment exists
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_ExistsByStudentAndCourse() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	enrollment := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	// Check before creating - should not exist
	exists, err := suite.repo.ExistsByStudentAndCourse(enrollment.StudentEmail, enrollment.CourseID)
	suite.NoError(err)
	suite.False(exists)

	// Create enrollment
	err = suite.repo.Create(enrollment)
	suite.NoError(err)

	// Check after creating - should exist
	exists, err = suite.repo.ExistsByStudentAndCourse(enrollment.StudentEmail, enrollment.CourseID)
	suite.NoError(err)
	suite.True(exists)
}

// TestEnrollmentRepository_Delete tests deleting an enrollment
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_Delete() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	enrollment := &models.Enrollment{
		ID:           uuid.New(),
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	err := suite.repo.Create(enrollment)
	suite.NoError(err)

	// Delete enrollment
	err = suite.repo.Delete(enrollment.ID)
	suite.NoError(err)

	// Verify deletion
	retrievedEnrollment, err := suite.repo.GetByStudentAndCourse(enrollment.StudentEmail, enrollment.CourseID)
	suite.Error(err)
	suite.Nil(retrievedEnrollment)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

// TestEnrollmentRepository_Delete_NotFound tests deleting a non-existent enrollment
func (suite *EnrollmentRepositoryTestSuite) TestEnrollmentRepository_Delete_NotFound() {
	nonExistentID := uuid.New()

	err := suite.repo.Delete(nonExistentID)

	// Should not return error for non-existent record (idempotent operation)
	suite.NoError(err)
}

// TestEnrollmentRepositoryTestSuite runs the enrollment repository test suite
func TestEnrollmentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EnrollmentRepositoryTestSuite))
}
