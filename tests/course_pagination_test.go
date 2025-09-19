package tests

import (
	"encoding/json"
	"net/http"
	"sonic-labs/course-enrollment-service/internal/models"
)

// TestCoursePagination tests pagination functionality
func (suite *IntegrationTestSuite) TestCoursePagination() {
	// Create test courses
	courses := []struct {
		title       string
		description string
		difficulty  string
	}{
		{"Go Programming", "Learn Go programming language", "Beginner"},
		{"Advanced Go", "Advanced Go concepts", "Advanced"},
		{"Python Basics", "Learn Python programming", "Beginner"},
		{"JavaScript Fundamentals", "Learn JavaScript", "Intermediate"},
		{"React Development", "Build React applications", "Intermediate"},
	}

	for _, course := range courses {
		suite.createTestCourse(course.title, course.description, course.difficulty)
	}

	// Test pagination - page 1, limit 2
	recorder := suite.makeRequest("GET", "/api/v1/courses?page=1&limit=2", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.CourseListResponse
	suite.parseResponse(recorder, &response)

	// Verify pagination metadata
	suite.Equal(1, response.Pagination.CurrentPage)
	suite.Equal(2, response.Pagination.Limit)
	suite.Equal(5, response.Pagination.TotalCount)
	suite.Equal(3, response.Pagination.TotalPages)
	suite.True(response.Pagination.HasNext)
	suite.False(response.Pagination.HasPrev)

	// Verify data
	suite.Len(response.Data, 2)

	// Test pagination - page 2, limit 2
	recorder = suite.makeRequest("GET", "/api/v1/courses?page=2&limit=2", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	suite.parseResponse(recorder, &response)

	// Verify pagination metadata
	suite.Equal(2, response.Pagination.CurrentPage)
	suite.True(response.Pagination.HasNext)
	suite.True(response.Pagination.HasPrev)

	// Verify data
	suite.Len(response.Data, 2)

	// Test pagination - last page
	recorder = suite.makeRequest("GET", "/api/v1/courses?page=3&limit=2", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	suite.parseResponse(recorder, &response)

	// Verify pagination metadata
	suite.Equal(3, response.Pagination.CurrentPage)
	suite.False(response.Pagination.HasNext)
	suite.True(response.Pagination.HasPrev)

	// Verify data (last page has 1 item)
	suite.Len(response.Data, 1)
}

// TestCourseSearch tests search functionality
func (suite *IntegrationTestSuite) TestCourseSearch() {
	// Create test courses
	courses := []struct {
		title       string
		description string
		difficulty  string
	}{
		{"Go Programming", "Learn Go programming language", "Beginner"},
		{"Advanced Go", "Advanced Go concepts and patterns", "Advanced"},
		{"Python Basics", "Learn Python programming", "Beginner"},
		{"JavaScript Fundamentals", "Learn JavaScript programming", "Intermediate"},
	}

	for _, course := range courses {
		suite.createTestCourse(course.title, course.description, course.difficulty)
	}

	// Test search by title
	recorder := suite.makeRequest("GET", "/api/v1/courses?search=Go", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.CourseListResponse
	suite.parseResponse(recorder, &response)

	// Should find 2 courses with "Go" in title
	suite.Equal(2, response.Pagination.TotalCount)
	suite.Len(response.Data, 2)

	// Test search by description
	recorder = suite.makeRequest("GET", "/api/v1/courses?search=programming", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	suite.parseResponse(recorder, &response)

	// Should find 3 courses with "programming" in description
	suite.Equal(3, response.Pagination.TotalCount)
	suite.Len(response.Data, 3)

	// Test search with no results
	recorder = suite.makeRequest("GET", "/api/v1/courses?search=nonexistent", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	suite.parseResponse(recorder, &response)

	// Should find no courses
	suite.Equal(0, response.Pagination.TotalCount)
	suite.Len(response.Data, 0)
}

// TestCourseDifficultyFilter tests difficulty filtering
func (suite *IntegrationTestSuite) TestCourseDifficultyFilter() {
	// Create test courses
	courses := []struct {
		title       string
		description string
		difficulty  string
	}{
		{"Go Basics", "Learn Go basics", "Beginner"},
		{"Go Advanced", "Advanced Go", "Advanced"},
		{"Python Basics", "Learn Python", "Beginner"},
		{"JavaScript Intermediate", "JS concepts", "Intermediate"},
		{"React Advanced", "Advanced React", "Advanced"},
	}

	for _, course := range courses {
		suite.createTestCourse(course.title, course.description, course.difficulty)
	}

	// Test filter by single difficulty
	recorder := suite.makeRequest("GET", "/api/v1/courses?difficulty=Beginner", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.CourseListResponse
	suite.parseResponse(recorder, &response)

	// Should find 2 beginner courses
	suite.Equal(2, response.Pagination.TotalCount)
	suite.Len(response.Data, 2)
	for _, course := range response.Data {
		suite.Equal("Beginner", course.Difficulty)
	}

	// Test filter by multiple difficulties
	recorder = suite.makeRequest("GET", "/api/v1/courses?difficulty=Beginner,Advanced", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	suite.parseResponse(recorder, &response)

	// Should find 4 courses (2 Beginner + 2 Advanced)
	suite.Equal(4, response.Pagination.TotalCount)
	suite.Len(response.Data, 4)
	for _, course := range response.Data {
		suite.True(course.Difficulty == "Beginner" || course.Difficulty == "Advanced")
	}
}

// TestCombinedFilters tests combination of search, difficulty filter, and pagination
func (suite *IntegrationTestSuite) TestCombinedFilters() {
	// Create test courses
	courses := []struct {
		title       string
		description string
		difficulty  string
	}{
		{"Go Programming Basics", "Learn Go programming", "Beginner"},
		{"Advanced Go Programming", "Advanced Go concepts", "Advanced"},
		{"Python Programming", "Learn Python programming", "Beginner"},
		{"Go Web Development", "Build web apps with Go", "Intermediate"},
	}

	for _, course := range courses {
		suite.createTestCourse(course.title, course.description, course.difficulty)
	}

	// Test search + difficulty filter + pagination
	url := "/api/v1/courses?search=Go&difficulty=Beginner,Intermediate&page=1&limit=2"
	recorder := suite.makeRequest("GET", url, nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	var response models.CourseListResponse
	suite.parseResponse(recorder, &response)

	// Should find 2 courses (Go Programming Basics + Go Web Development)
	suite.Equal(2, response.Pagination.TotalCount)
	suite.Len(response.Data, 2)

	// Verify all courses contain "Go" and have correct difficulty
	for _, course := range response.Data {
		suite.Contains(course.Title, "Go")
		suite.True(course.Difficulty == "Beginner" || course.Difficulty == "Intermediate")
	}
}

// TestBackwardCompatibility tests that existing API calls still work
func (suite *IntegrationTestSuite) TestBackwardCompatibility() {
	// Create test courses
	suite.createTestCourse("Course 1", "Description 1", "Beginner")
	suite.createTestCourse("Course 2", "Description 2", "Intermediate")

	// Test old API call without parameters (should return simple array)
	recorder := suite.makeRequest("GET", "/api/v1/courses", nil, suite.getAuthHeaders())
	suite.Equal(http.StatusOK, recorder.Code)

	// Should return simple array, not paginated response
	var courses []models.CourseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &courses)
	suite.NoError(err)
	suite.Len(courses, 2)

	// Verify it's not a paginated response
	var paginatedResponse models.CourseListResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &paginatedResponse)
	// This should fail because it's not a paginated response
	suite.Error(err)
}
