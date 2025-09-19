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

// CourseRepositoryTestSuite defines the test suite for course repository tests
type CourseRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo CourseRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *CourseRepositoryTestSuite) SetupSuite() {
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

	// Initialize repository
	suite.repo = NewCourseRepository(suite.db)
}

// TearDownSuite runs once after all tests in the suite
func (suite *CourseRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// SetupTest runs before each test
func (suite *CourseRepositoryTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.db.Exec("DELETE FROM courses")
}

// TestCourseRepository_Create tests creating a course
func (suite *CourseRepositoryTestSuite) TestCourseRepository_Create() {
	course := &models.Course{
		ID:          uuid.New(),
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	err := suite.repo.Create(course)

	suite.NoError(err)
	suite.NotEqual(uuid.Nil, course.ID)

	// Verify course was created in database
	var dbCourse models.Course
	err = suite.db.First(&dbCourse, "id = ?", course.ID.String()).Error
	suite.NoError(err)
	suite.Equal(course.Title, dbCourse.Title)
	suite.Equal(course.Description, dbCourse.Description)
	suite.Equal(course.Difficulty, dbCourse.Difficulty)
}

// TestCourseRepository_GetAll tests retrieving all courses
func (suite *CourseRepositoryTestSuite) TestCourseRepository_GetAll() {
	// Create test courses
	course1 := &models.Course{
		ID:          uuid.New(),
		Title:       "Course 1",
		Description: "Description 1",
		Difficulty:  "Beginner",
	}
	course2 := &models.Course{
		ID:          uuid.New(),
		Title:       "Course 2",
		Description: "Description 2",
		Difficulty:  "Intermediate",
	}

	err := suite.repo.Create(course1)
	suite.NoError(err)
	err = suite.repo.Create(course2)
	suite.NoError(err)

	// Get all courses
	courses, err := suite.repo.GetAll()

	suite.NoError(err)
	suite.Len(courses, 2)

	// Verify course data (order might vary)
	courseMap := make(map[string]models.Course)
	for _, course := range courses {
		courseMap[course.Title] = course
	}

	suite.Contains(courseMap, "Course 1")
	suite.Contains(courseMap, "Course 2")
	suite.Equal("Description 1", courseMap["Course 1"].Description)
	suite.Equal("Beginner", courseMap["Course 1"].Difficulty)
	suite.Equal("Description 2", courseMap["Course 2"].Description)
	suite.Equal("Intermediate", courseMap["Course 2"].Difficulty)
}

// TestCourseRepository_GetAll_Empty tests retrieving all courses when none exist
func (suite *CourseRepositoryTestSuite) TestCourseRepository_GetAll_Empty() {
	courses, err := suite.repo.GetAll()

	suite.NoError(err)
	suite.Len(courses, 0)
}

// TestCourseRepository_GetByID tests retrieving a course by ID
func (suite *CourseRepositoryTestSuite) TestCourseRepository_GetByID() {
	// Create test course
	course := &models.Course{
		ID:          uuid.New(),
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	err := suite.repo.Create(course)
	suite.NoError(err)

	// Get course by ID
	retrievedCourse, err := suite.repo.GetByID(course.ID)

	suite.NoError(err)
	suite.NotNil(retrievedCourse)
	suite.Equal(course.ID, retrievedCourse.ID)
	suite.Equal(course.Title, retrievedCourse.Title)
	suite.Equal(course.Description, retrievedCourse.Description)
	suite.Equal(course.Difficulty, retrievedCourse.Difficulty)
}

// TestCourseRepository_GetByID_NotFound tests retrieving a non-existent course
func (suite *CourseRepositoryTestSuite) TestCourseRepository_GetByID_NotFound() {
	nonExistentID := uuid.New()

	retrievedCourse, err := suite.repo.GetByID(nonExistentID)

	suite.Error(err)
	suite.Nil(retrievedCourse)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

// TestCourseRepository_Update tests updating a course
func (suite *CourseRepositoryTestSuite) TestCourseRepository_Update() {
	// Create test course
	course := &models.Course{
		ID:          uuid.New(),
		Title:       "Original Title",
		Description: "Original Description",
		Difficulty:  "Beginner",
	}

	err := suite.repo.Create(course)
	suite.NoError(err)

	// Update course
	course.Title = "Updated Title"
	course.Description = "Updated Description"
	course.Difficulty = "Advanced"

	err = suite.repo.Update(course)
	suite.NoError(err)

	// Verify update
	retrievedCourse, err := suite.repo.GetByID(course.ID)
	suite.NoError(err)
	suite.Equal("Updated Title", retrievedCourse.Title)
	suite.Equal("Updated Description", retrievedCourse.Description)
	suite.Equal("Advanced", retrievedCourse.Difficulty)
}

// TestCourseRepository_Delete tests deleting a course
func (suite *CourseRepositoryTestSuite) TestCourseRepository_Delete() {
	// Create test course
	course := &models.Course{
		ID:          uuid.New(),
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	err := suite.repo.Create(course)
	suite.NoError(err)

	// Delete course
	err = suite.repo.Delete(course.ID)
	suite.NoError(err)

	// Verify deletion
	retrievedCourse, err := suite.repo.GetByID(course.ID)
	suite.Error(err)
	suite.Nil(retrievedCourse)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

// TestCourseRepository_Delete_NotFound tests deleting a non-existent course
func (suite *CourseRepositoryTestSuite) TestCourseRepository_Delete_NotFound() {
	nonExistentID := uuid.New()

	err := suite.repo.Delete(nonExistentID)

	// Should not return error for non-existent record (idempotent operation)
	suite.NoError(err)
}

// TestCourseRepositoryTestSuite runs the course repository test suite
func TestCourseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourseRepositoryTestSuite))
}
