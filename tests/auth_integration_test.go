package tests

import (
	"bytes"
	"encoding/json"
	"net/http"

	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
)

// TestAuthLogin tests the authentication login endpoint
func (suite *IntegrationTestSuite) TestAuthLogin() {
	// Test successful login
	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "admin!dev",
	}

	resp := suite.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
	suite.Equal(http.StatusOK, resp.Code)

	var loginResp models.LoginResponse
	err := json.NewDecoder(resp.Body).Decode(&loginResp)
	suite.NoError(err)
	suite.NotEmpty(loginResp.Token)
	suite.Equal("admin", loginResp.User.Username)
	suite.Equal("admin", loginResp.User.Role)
	suite.NotEqual(uuid.Nil, loginResp.User.ID)
}

// TestAuthLoginInvalidCredentials tests login with invalid credentials
func (suite *IntegrationTestSuite) TestAuthLoginInvalidCredentials() {
	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	resp := suite.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
	suite.Equal(http.StatusUnauthorized, resp.Code)

	var errorResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	suite.NoError(err)
	suite.Contains(errorResp["message"], "Invalid username or password")
}

// TestAuthProfile tests the user profile endpoint
func (suite *IntegrationTestSuite) TestAuthProfile() {
	// First login to get token
	token := suite.getAuthToken()

	// Test getting profile with valid token
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp := suite.makeRequest("GET", "/api/v1/auth/profile", nil, headers)
	suite.Equal(http.StatusOK, resp.Code)

	var userResp models.UserResponse
	err := json.NewDecoder(resp.Body).Decode(&userResp)
	suite.NoError(err)
	suite.Equal("admin", userResp.Username)
	suite.Equal("admin", userResp.Role)
	suite.NotEqual(uuid.Nil, userResp.ID)
}

// TestAuthProfileWithoutToken tests profile endpoint without token
func (suite *IntegrationTestSuite) TestAuthProfileWithoutToken() {
	resp := suite.makeRequest("GET", "/api/v1/auth/profile", nil, nil)
	suite.Equal(http.StatusUnauthorized, resp.Code)

	var errorResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	suite.NoError(err)
	suite.Contains(errorResp["message"], "Authorization header is required")
}

// TestAuthProfileWithInvalidToken tests profile endpoint with invalid token
func (suite *IntegrationTestSuite) TestAuthProfileWithInvalidToken() {
	headers := map[string]string{
		"Authorization": "Bearer invalid-token",
	}

	resp := suite.makeRequest("GET", "/api/v1/auth/profile", nil, headers)
	suite.Equal(http.StatusUnauthorized, resp.Code)

	var errorResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	suite.NoError(err)
	suite.Contains(errorResp["message"], "JWT token is invalid or expired")
}

// TestProtectedCourseCreation tests creating a course with authentication
func (suite *IntegrationTestSuite) TestProtectedCourseCreation() {
	token := suite.getAuthToken()

	courseReq := models.CourseRequest{
		Title:       "Protected Test Course",
		Description: "This course was created with authentication",
		Difficulty:  "Beginner",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp := suite.makeRequest("POST", "/api/v1/courses", courseReq, headers)
	suite.Equal(http.StatusCreated, resp.Code)

	var courseResp models.CourseResponse
	err := json.NewDecoder(resp.Body).Decode(&courseResp)
	suite.NoError(err)
	suite.Equal(courseReq.Title, courseResp.Title)
	suite.Equal(courseReq.Description, courseResp.Description)
	suite.Equal(courseReq.Difficulty, courseResp.Difficulty)
	suite.NotEqual(uuid.Nil, courseResp.ID)
}

// TestProtectedEnrollmentCreation tests creating an enrollment with authentication
func (suite *IntegrationTestSuite) TestProtectedEnrollmentCreation() {
	token := suite.getAuthToken()

	// First create a course
	course := suite.createTestCourse("Test Course for Enrollment", "Test Description", "Beginner")

	enrollmentReq := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp := suite.makeRequest("POST", "/api/v1/enrollments", enrollmentReq, headers)
	suite.Equal(http.StatusCreated, resp.Code)

	var enrollmentResp models.EnrollmentResponse
	err := json.NewDecoder(resp.Body).Decode(&enrollmentResp)
	suite.NoError(err)
	suite.Equal(enrollmentReq.StudentEmail, enrollmentResp.StudentEmail)
	suite.Equal(enrollmentReq.CourseID, enrollmentResp.CourseID)
	suite.NotEqual(uuid.Nil, enrollmentResp.ID)
}

// TestUnauthorizedAccess tests accessing protected endpoints without authentication
func (suite *IntegrationTestSuite) TestUnauthorizedAccess() {
	// Test creating course without token
	courseReq := models.CourseRequest{
		Title:       "Unauthorized Course",
		Description: "This should fail",
		Difficulty:  "Beginner",
	}

	resp := suite.makeRequest("POST", "/api/v1/courses", courseReq, nil)
	suite.Equal(http.StatusUnauthorized, resp.Code)

	// Test creating enrollment without token
	enrollmentReq := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     uuid.New(),
	}

	resp = suite.makeRequest("POST", "/api/v1/enrollments", enrollmentReq, nil)
	suite.Equal(http.StatusUnauthorized, resp.Code)
}

// getAuthToken is a helper method to get a valid authentication token
func (suite *IntegrationTestSuite) getAuthToken() string {
	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "admin!dev",
	}

	reqBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp := suite.makeHTTPRequest(req)

	var loginResp models.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)

	return loginResp.Token
}
