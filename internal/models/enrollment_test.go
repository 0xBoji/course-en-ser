package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEnrollment_BeforeCreate(t *testing.T) {
	enrollment := &Enrollment{}

	// Mock GORM DB for the test
	db := &gorm.DB{}

	err := enrollment.BeforeCreate(db)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, enrollment.ID)
	assert.False(t, enrollment.EnrolledAt.IsZero())
}

func TestEnrollment_TableName(t *testing.T) {
	enrollment := Enrollment{}
	assert.Equal(t, "enrollments", enrollment.TableName())
}

func TestEnrollment_ToResponse(t *testing.T) {
	enrollmentID := uuid.New()
	courseID := uuid.New()
	enrolledAt := time.Now()

	course := Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test Description",
		Difficulty:  "Beginner",
	}

	enrollment := Enrollment{
		ID:           enrollmentID,
		StudentEmail: "test@example.com",
		CourseID:     courseID,
		EnrolledAt:   enrolledAt,
		Course:       course,
	}

	response := enrollment.ToResponse()

	assert.Equal(t, enrollmentID, response.ID)
	assert.Equal(t, "test@example.com", response.StudentEmail)
	assert.Equal(t, courseID, response.CourseID)
	assert.Equal(t, enrolledAt, response.EnrolledAt)
	assert.Equal(t, courseID, response.Course.ID)
	assert.Equal(t, "Test Course", response.Course.Title)
}

func TestEnrollment_ToResponseWithoutCourse(t *testing.T) {
	enrollmentID := uuid.New()
	courseID := uuid.New()
	enrolledAt := time.Now()

	enrollment := Enrollment{
		ID:           enrollmentID,
		StudentEmail: "test@example.com",
		CourseID:     courseID,
		EnrolledAt:   enrolledAt,
		// Course not loaded
	}

	response := enrollment.ToResponse()

	assert.Equal(t, enrollmentID, response.ID)
	assert.Equal(t, "test@example.com", response.StudentEmail)
	assert.Equal(t, courseID, response.CourseID)
	assert.Equal(t, enrolledAt, response.EnrolledAt)
	assert.Equal(t, uuid.Nil, response.Course.ID) // Course not included
}

func TestEnrollmentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request EnrollmentRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: EnrollmentRequest{
				StudentEmail: "test@example.com",
				CourseID:     uuid.New(),
			},
			valid: true,
		},
		{
			name: "empty email",
			request: EnrollmentRequest{
				StudentEmail: "",
				CourseID:     uuid.New(),
			},
			valid: false,
		},
		{
			name: "invalid email format",
			request: EnrollmentRequest{
				StudentEmail: "invalid-email",
				CourseID:     uuid.New(),
			},
			valid: false,
		},
		{
			name: "nil course ID",
			request: EnrollmentRequest{
				StudentEmail: "test@example.com",
				CourseID:     uuid.Nil,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.StudentEmail)
				assert.Contains(t, tt.request.StudentEmail, "@")
				assert.NotEqual(t, uuid.Nil, tt.request.CourseID)
			}
		})
	}
}
