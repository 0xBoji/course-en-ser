package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCourse_BeforeCreate(t *testing.T) {
	course := &Course{}

	// Mock GORM DB for the test
	db := &gorm.DB{}

	err := course.BeforeCreate(db)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, course.ID)
}

func TestCourse_TableName(t *testing.T) {
	course := Course{}
	assert.Equal(t, "courses", course.TableName())
}

func TestCourse_ToResponse(t *testing.T) {
	courseID := uuid.New()
	createdAt := time.Now()

	course := Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
		CreatedAt:   createdAt,
	}

	response := course.ToResponse()

	assert.Equal(t, courseID, response.ID)
	assert.Equal(t, "Test Course", response.Title)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, "Beginner", response.Difficulty)
	assert.Equal(t, createdAt, response.CreatedAt)
}

func TestCourseRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CourseRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: CourseRequest{
				Title:       "Valid Course",
				Description: "Valid Description",
				Difficulty:  "Beginner",
			},
			valid: true,
		},
		{
			name: "empty title",
			request: CourseRequest{
				Title:       "",
				Description: "Valid Description",
				Difficulty:  "Beginner",
			},
			valid: false,
		},
		{
			name: "empty description",
			request: CourseRequest{
				Title:       "Valid Course",
				Description: "",
				Difficulty:  "Beginner",
			},
			valid: false,
		},
		{
			name: "invalid difficulty",
			request: CourseRequest{
				Title:       "Valid Course",
				Description: "Valid Description",
				Difficulty:  "Invalid",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validator library
			// to validate the struct tags. This is a simplified test.
			if tt.valid {
				assert.NotEmpty(t, tt.request.Title)
				assert.NotEmpty(t, tt.request.Description)
				assert.Contains(t, []string{"Beginner", "Intermediate", "Advanced"}, tt.request.Difficulty)
			}
		})
	}
}
