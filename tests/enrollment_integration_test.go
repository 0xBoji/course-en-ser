package tests

import (
	"fmt"
	"net/http"

	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
)

// TestEnrollStudent tests the POST /api/v1/enrollments endpoint - happy path
func (suite *IntegrationTestSuite) TestEnrollStudent() {
	// Create test course
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	// Prepare enrollment request
	enrollReq := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	// Make request with authentication
	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/enrollments", enrollReq, headers)

	// Assert response
	suite.Equal(http.StatusCreated, recorder.Code)

	var enrollment models.EnrollmentResponse
	suite.parseResponse(recorder, &enrollment)

	suite.Equal(enrollReq.StudentEmail, enrollment.StudentEmail)
	suite.Equal(enrollReq.CourseID, enrollment.CourseID)
	suite.NotEqual(uuid.Nil, enrollment.ID)
	suite.False(enrollment.EnrolledAt.IsZero())

	// Verify enrollment was actually created in database
	var dbEnrollment models.Enrollment
	err := suite.db.First(&dbEnrollment, "id = ?", enrollment.ID).Error
	suite.NoError(err)
	suite.Equal(enrollReq.StudentEmail, dbEnrollment.StudentEmail)
	suite.Equal(enrollReq.CourseID, dbEnrollment.CourseID)
}

// TestEnrollStudentDuplicate tests POST /api/v1/enrollments with duplicate enrollment (409 Conflict)
func (suite *IntegrationTestSuite) TestEnrollStudentDuplicate() {
	// Create test course
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	// Create initial enrollment
	enrollReq := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     course.ID,
	}

	// First enrollment should succeed
	headers := suite.getAuthHeaders()
	recorder1 := suite.makeRequest("POST", "/api/v1/enrollments", enrollReq, headers)
	suite.Equal(http.StatusCreated, recorder1.Code)

	// Second enrollment should fail with conflict
	recorder2 := suite.makeRequest("POST", "/api/v1/enrollments", enrollReq, headers)
	suite.assertErrorResponse(recorder2, http.StatusConflict, "Student is already enrolled")
}

// TestEnrollStudentValidationErrors tests POST /api/v1/enrollments with validation errors
func (suite *IntegrationTestSuite) TestEnrollStudentValidationErrors() {
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	testCases := []struct {
		name           string
		request        models.EnrollmentRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing email",
			request: models.EnrollmentRequest{
				StudentEmail: "",
				CourseID:     course.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Student email is required",
		},
		{
			name: "invalid email format",
			request: models.EnrollmentRequest{
				StudentEmail: "invalid-email",
				CourseID:     course.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "missing course ID",
			request: models.EnrollmentRequest{
				StudentEmail: "student@example.com",
				CourseID:     uuid.Nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Course ID is required",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			headers := suite.getAuthHeaders()
			recorder := suite.makeRequest("POST", "/api/v1/enrollments", tc.request, headers)
			suite.assertErrorResponse(recorder, tc.expectedStatus, tc.expectedError)
		})
	}
}

// TestEnrollStudentCourseNotFound tests POST /api/v1/enrollments with non-existent course
func (suite *IntegrationTestSuite) TestEnrollStudentCourseNotFound() {
	// Use a random UUID that doesn't exist
	nonExistentCourseID := uuid.New()

	enrollReq := models.EnrollmentRequest{
		StudentEmail: "student@example.com",
		CourseID:     nonExistentCourseID,
	}

	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/enrollments", enrollReq, headers)
	suite.assertErrorResponse(recorder, http.StatusBadRequest, "The specified course does not exist")
}

// TestEnrollStudentInvalidJSON tests POST /api/v1/enrollments with invalid JSON
func (suite *IntegrationTestSuite) TestEnrollStudentInvalidJSON() {
	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/enrollments", "invalid json", headers)
	suite.assertErrorResponse(recorder, http.StatusBadRequest, "")
}

// TestGetStudentEnrollments tests the GET /api/v1/students/:email/enrollments endpoint
func (suite *IntegrationTestSuite) TestGetStudentEnrollments() {
	// Create test courses
	course1 := suite.createTestCourse("Course 1", "Description 1", "Beginner")
	course2 := suite.createTestCourse("Course 2", "Description 2", "Intermediate")

	studentEmail := "student@example.com"

	// Create enrollments
	enrollment1 := &models.Enrollment{
		StudentEmail: studentEmail,
		CourseID:     course1.ID,
	}
	enrollment2 := &models.Enrollment{
		StudentEmail: studentEmail,
		CourseID:     course2.ID,
	}

	err := suite.db.Create(enrollment1).Error
	suite.NoError(err)
	err = suite.db.Create(enrollment2).Error
	suite.NoError(err)

	// Make request
	url := fmt.Sprintf("/api/v1/students/%s/enrollments", studentEmail)
	recorder := suite.makeRequest("GET", url, nil, nil)

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.StudentEnrollmentsResponse
	suite.parseResponse(recorder, &response)

	suite.Equal(studentEmail, response.StudentEmail)
	suite.Equal(2, response.Total)
	suite.Len(response.Enrollments, 2)

	// Verify enrollment data (order might vary)
	enrollmentMap := make(map[string]models.EnrollmentResponse)
	for _, enrollment := range response.Enrollments {
		enrollmentMap[enrollment.Course.Title] = enrollment
	}

	suite.Contains(enrollmentMap, "Course 1")
	suite.Contains(enrollmentMap, "Course 2")
	suite.Equal(course1.ID, enrollmentMap["Course 1"].CourseID)
	suite.Equal(course2.ID, enrollmentMap["Course 2"].CourseID)
}

// TestGetStudentEnrollmentsEmpty tests GET /api/v1/students/:email/enrollments with no enrollments
func (suite *IntegrationTestSuite) TestGetStudentEnrollmentsEmpty() {
	studentEmail := "student@example.com"

	// Make request without creating any enrollments
	url := fmt.Sprintf("/api/v1/students/%s/enrollments", studentEmail)
	recorder := suite.makeRequest("GET", url, nil, nil)

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.StudentEnrollmentsResponse
	suite.parseResponse(recorder, &response)

	suite.Equal(studentEmail, response.StudentEmail)
	suite.Equal(0, response.Total)
	suite.Len(response.Enrollments, 0)
}

// TestGetStudentEnrollmentsInvalidEmail tests GET /api/v1/students/:email/enrollments with invalid email
func (suite *IntegrationTestSuite) TestGetStudentEnrollmentsInvalidEmail() {
	invalidEmail := "invalid-email"

	url := fmt.Sprintf("/api/v1/students/%s/enrollments", invalidEmail)
	recorder := suite.makeRequest("GET", url, nil, nil)

	suite.assertErrorResponse(recorder, http.StatusBadRequest, "Invalid email format")
}

// TestGetStudentEnrollmentsWithCourseDetails tests that course details are included in enrollment response
func (suite *IntegrationTestSuite) TestGetStudentEnrollmentsWithCourseDetails() {
	// Create test course
	course := suite.createTestCourse("Detailed Course", "Detailed Description", "Advanced")

	studentEmail := "student@example.com"

	// Create enrollment
	enrollment := &models.Enrollment{
		StudentEmail: studentEmail,
		CourseID:     course.ID,
	}

	err := suite.db.Create(enrollment).Error
	suite.NoError(err)

	// Make request
	url := fmt.Sprintf("/api/v1/students/%s/enrollments", studentEmail)
	recorder := suite.makeRequest("GET", url, nil, nil)

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.StudentEnrollmentsResponse
	suite.parseResponse(recorder, &response)

	suite.Len(response.Enrollments, 1)

	enrollmentResp := response.Enrollments[0]
	suite.Equal(course.ID, enrollmentResp.Course.ID)
	suite.Equal(course.Title, enrollmentResp.Course.Title)
	suite.Equal(course.Description, enrollmentResp.Course.Description)
	suite.Equal(course.Difficulty, enrollmentResp.Course.Difficulty)
}
