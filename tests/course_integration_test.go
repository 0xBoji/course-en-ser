package tests

import (
	"fmt"
	"net/http"

	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/google/uuid"
)

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// TestGetAllCourses tests the GET /api/v1/courses endpoint
func (suite *IntegrationTestSuite) TestGetAllCourses() {
	// Create test courses
	suite.createTestCourse("Course 1", "Description 1", "Beginner")
	suite.createTestCourse("Course 2", "Description 2", "Intermediate")

	// Make request with authentication
	recorder := suite.makeRequest("GET", "/api/v1/courses", nil, suite.getAuthHeaders())

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var courses []models.CourseResponse
	suite.parseResponse(recorder, &courses)

	suite.Len(courses, 2)

	// Verify course data (order might vary, so check both possibilities)
	courseMap := make(map[string]models.CourseResponse)
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

// TestGetAllCoursesEmpty tests GET /api/v1/courses with no courses
func (suite *IntegrationTestSuite) TestGetAllCoursesEmpty() {
	// Make request without creating any courses
	recorder := suite.makeRequest("GET", "/api/v1/courses", nil, suite.getAuthHeaders())

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var courses []models.CourseResponse
	suite.parseResponse(recorder, &courses)

	suite.Len(courses, 0)
}

// TestCreateCourse tests the POST /api/v1/courses endpoint - happy path
func (suite *IntegrationTestSuite) TestCreateCourse() {
	// Prepare request
	courseReq := models.CourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	// Make request with authentication
	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/courses", courseReq, headers)

	// Assert response
	suite.Equal(http.StatusCreated, recorder.Code)

	var course models.CourseResponse
	suite.parseResponse(recorder, &course)

	suite.Equal(courseReq.Title, course.Title)
	suite.Equal(courseReq.Description, course.Description)
	suite.Equal(courseReq.Difficulty, course.Difficulty)
	suite.Nil(course.ImageURL) // Should be nil when not provided
	suite.NotEqual(uuid.Nil, course.ID)
	suite.False(course.CreatedAt.IsZero())

	// Verify course was actually created in database
	var dbCourse models.Course
	err := suite.db.First(&dbCourse, "id = ?", course.ID).Error
	suite.NoError(err)
	suite.Equal(courseReq.Title, dbCourse.Title)
}

// TestCreateCourseWithImage tests the POST /api/v1/courses endpoint with image URL
func (suite *IntegrationTestSuite) TestCreateCourseWithImage() {
	imageURL := "https://example.com/course-image.jpg"
	
	// Prepare request
	courseReq := models.CourseRequest{
		Title:       "Test Course with Image",
		Description: "Test Description with Image",
		Difficulty:  "Intermediate",
		ImageURL:    &imageURL,
	}

	// Make request with authentication
	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/courses", courseReq, headers)

	// Assert response
	suite.Equal(http.StatusCreated, recorder.Code)

	var course models.CourseResponse
	suite.parseResponse(recorder, &course)

	suite.Equal(courseReq.Title, course.Title)
	suite.Equal(courseReq.Description, course.Description)
	suite.Equal(courseReq.Difficulty, course.Difficulty)
	suite.NotNil(course.ImageURL)
	suite.Equal(imageURL, *course.ImageURL)
	suite.NotEqual(uuid.Nil, course.ID)
	suite.False(course.CreatedAt.IsZero())

	// Verify course was actually created in database
	var dbCourse models.Course
	err := suite.db.First(&dbCourse, "id = ?", course.ID).Error
	suite.NoError(err)
	suite.Equal(courseReq.Title, dbCourse.Title)
	suite.NotNil(dbCourse.ImageURL)
	suite.Equal(imageURL, *dbCourse.ImageURL)
}

// TestCreateCourseValidationErrors tests POST /api/v1/courses with validation errors
func (suite *IntegrationTestSuite) TestCreateCourseValidationErrors() {
	testCases := []struct {
		name           string
		request        models.CourseRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing title",
			request: models.CourseRequest{
				Title:       "",
				Description: "Valid Description",
				Difficulty:  "Beginner",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Title is required",
		},
		{
			name: "missing description",
			request: models.CourseRequest{
				Title:       "Valid Title",
				Description: "",
				Difficulty:  "Beginner",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Description is required",
		},
		{
			name: "invalid difficulty",
			request: models.CourseRequest{
				Title:       "Valid Title",
				Description: "Valid Description",
				Difficulty:  "Invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Difficulty must be one of",
		},
		{
			name: "missing difficulty",
			request: models.CourseRequest{
				Title:       "Valid Title",
				Description: "Valid Description",
				Difficulty:  "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Difficulty must be one of",
		},
		{
			name: "invalid image URL",
			request: models.CourseRequest{
				Title:       "Valid Title",
				Description: "Valid Description",
				Difficulty:  "Beginner",
				ImageURL:    stringPtr("not-a-valid-url"),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Image URL must be a valid URL",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			headers := suite.getAuthHeaders()
			recorder := suite.makeRequest("POST", "/api/v1/courses", tc.request, headers)
			suite.assertErrorResponse(recorder, tc.expectedStatus, tc.expectedError)
		})
	}
}

// TestCreateCourseInvalidJSON tests POST /api/v1/courses with invalid JSON
func (suite *IntegrationTestSuite) TestCreateCourseInvalidJSON() {
	headers := suite.getAuthHeaders()
	recorder := suite.makeRequest("POST", "/api/v1/courses", "invalid json", headers)
	suite.assertErrorResponse(recorder, http.StatusBadRequest, "")
}

// TestGetCourseByID tests the GET /api/v1/courses/:id endpoint
func (suite *IntegrationTestSuite) TestGetCourseByID() {
	// Create test course
	course := suite.createTestCourse("Test Course", "Test Description", "Beginner")

	// Make request
	url := fmt.Sprintf("/api/v1/courses/%s", course.ID.String())
	recorder := suite.makeRequest("GET", url, nil, suite.getAuthHeaders())

	// Assert response
	suite.Equal(http.StatusOK, recorder.Code)

	var courseResp models.CourseResponse
	suite.parseResponse(recorder, &courseResp)

	suite.Equal(course.ID, courseResp.ID)
	suite.Equal(course.Title, courseResp.Title)
	suite.Equal(course.Description, courseResp.Description)
	suite.Equal(course.Difficulty, courseResp.Difficulty)
}

// TestGetCourseByIDNotFound tests GET /api/v1/courses/:id with non-existent ID
func (suite *IntegrationTestSuite) TestGetCourseByIDNotFound() {
	// Use a random UUID that doesn't exist
	nonExistentID := uuid.New()
	url := fmt.Sprintf("/api/v1/courses/%s", nonExistentID.String())

	recorder := suite.makeRequest("GET", url, nil, suite.getAuthHeaders())
	suite.assertErrorResponse(recorder, http.StatusNotFound, "The requested course does not exist")
}

// TestGetCourseByIDInvalidUUID tests GET /api/v1/courses/:id with invalid UUID
func (suite *IntegrationTestSuite) TestGetCourseByIDInvalidUUID() {
	url := "/api/v1/courses/invalid-uuid"
	recorder := suite.makeRequest("GET", url, nil, suite.getAuthHeaders())
	suite.assertErrorResponse(recorder, http.StatusBadRequest, "")
}
