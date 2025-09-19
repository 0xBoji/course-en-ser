package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()" example:"123e4567-e89b-12d3-a456-426614174000"`
	StudentEmail string    `json:"student_email" gorm:"not null;size:255;index:idx_student_course,unique" validate:"required,email" example:"student@example.com"`
	CourseID     uuid.UUID `json:"course_id" gorm:"type:uuid;not null;index:idx_student_course,unique" example:"123e4567-e89b-12d3-a456-426614174000"`
	EnrolledAt   time.Time `json:"enrolled_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`

	// Relationships
	Course Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.EnrolledAt.IsZero() {
		e.EnrolledAt = time.Now()
	}
	return nil
}

// TableName returns the table name for Enrollment model
func (Enrollment) TableName() string {
	return "enrollments"
}

// EnrollmentRequest represents the request payload for creating an enrollment
type EnrollmentRequest struct {
	StudentEmail string    `json:"student_email" validate:"required,email" example:"student@example.com"`
	CourseID     uuid.UUID `json:"course_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// EnrollmentResponse represents the response payload for enrollment operations
type EnrollmentResponse struct {
	ID           uuid.UUID      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StudentEmail string         `json:"student_email" example:"student@example.com"`
	CourseID     uuid.UUID      `json:"course_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	EnrolledAt   time.Time      `json:"enrolled_at" example:"2023-01-01T00:00:00Z"`
	Course       CourseResponse `json:"course,omitempty"`
}

// ToResponse converts Enrollment model to EnrollmentResponse
func (e *Enrollment) ToResponse() EnrollmentResponse {
	response := EnrollmentResponse{
		ID:           e.ID,
		StudentEmail: e.StudentEmail,
		CourseID:     e.CourseID,
		EnrolledAt:   e.EnrolledAt,
	}

	// Include course information if loaded
	if e.Course.ID != uuid.Nil {
		response.Course = e.Course.ToResponse()
	}

	return response
}

// StudentEnrollmentsResponse represents the response for student's enrollments
type StudentEnrollmentsResponse struct {
	StudentEmail string               `json:"student_email" example:"student@example.com"`
	Enrollments  []EnrollmentResponse `json:"enrollments"`
	Total        int                  `json:"total" example:"3"`
}
