package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"sonic-labs/course-enrollment-service/internal/auth"
	"sonic-labs/course-enrollment-service/internal/config"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// IntegrationTestSuite defines the test suite for integration tests
type IntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	cfg    *config.Config
}

// SetupSuite runs once before all tests in the suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test configuration
	suite.cfg = &config.Config{
		Port:      "8080",
		JWTSecret: "test-jwt-secret-for-integration-tests",
	}

	// Initialize JWT secret
	auth.SetJWTSecret(suite.cfg.JWTSecret)

	// Initialize in-memory SQLite database for testing
	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Run migrations with SQLite-compatible schema
	err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS courses (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			difficulty TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		log.Fatalf("Failed to create courses table: %v", err)
	}

	err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS enrollments (
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
	if err != nil {
		log.Fatalf("Failed to create enrollments table: %v", err)
	}

	err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'admin',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create admin user for testing
	// Password is hashed using bcrypt for 'admin!dev'
	err = suite.db.Exec(`
		INSERT OR IGNORE INTO users (id, username, password, role) 
		VALUES ('12345678-1234-1234-1234-123456789012', 'admin', '$2a$10$V6C81VGFyKg/sRc1JOw8cOs7dV/3StzYs5NUZaYvDFcEEKW0Tlika', 'admin')
	`).Error
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	// Setup router
	suite.router = router.Setup(suite.db, suite.cfg)
}

// TearDownSuite runs once after all tests in the suite
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up database connection
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *IntegrationTestSuite) TearDownTest() {
	// Clean up test data after each test
	suite.cleanupTestData()
}

// cleanupTestData removes all test data from the database
func (suite *IntegrationTestSuite) cleanupTestData() {
	// Delete in order to respect foreign key constraints
	suite.db.Exec("DELETE FROM enrollments")
	suite.db.Exec("DELETE FROM courses")
	// Don't delete users as we need admin user for tests
}

// makeRequest is a helper function to make HTTP requests to the test server
func (suite *IntegrationTestSuite) makeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		suite.Require().NoError(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	suite.Require().NoError(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	return recorder
}

// makeHTTPRequest is a helper function for making raw HTTP requests
func (suite *IntegrationTestSuite) makeHTTPRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)
	return recorder
}

// createTestCourse is a helper function to create a test course
func (suite *IntegrationTestSuite) createTestCourse(title, description, difficulty string) *models.Course {
	course := &models.Course{
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
	}

	err := suite.db.Create(course).Error
	suite.Require().NoError(err)

	return course
}

// parseResponse is a helper function to parse JSON response
func (suite *IntegrationTestSuite) parseResponse(recorder *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(recorder.Body.Bytes(), target)
	suite.Require().NoError(err)
}

// getAuthHeaders is a helper method to get authentication headers
func (suite *IntegrationTestSuite) getAuthHeaders() map[string]string {
	token := suite.getAuthToken()
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// assertErrorResponse is a helper function to assert error response format
func (suite *IntegrationTestSuite) assertErrorResponse(recorder *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	suite.Equal(expectedStatus, recorder.Code)

	var errorResp map[string]interface{}
	suite.parseResponse(recorder, &errorResp)

	suite.Contains(errorResp, "error")
	if expectedError != "" {
		// Check the message field for the expected text
		if errorResp["message"] != nil {
			messageField := fmt.Sprintf("%v", errorResp["message"])
			suite.Contains(messageField, expectedError)
		}
	}
}

// TestIntegrationTestSuite runs the integration test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
