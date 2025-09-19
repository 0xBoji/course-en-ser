package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Course represents a course in the system
type Course struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string    `json:"title" gorm:"not null;size:255" validate:"required,min=1,max=255" example:"Introduction to Go Programming"`
	Description string    `json:"description" gorm:"not null;type:text" validate:"required,min=1" example:"Learn the fundamentals of Go programming language"`
	Difficulty  string    `json:"difficulty" gorm:"not null;size:50" validate:"required,oneof=Beginner Intermediate Advanced" example:"Beginner"`
	ImageURL    *string   `json:"image_url,omitempty" gorm:"size:500" validate:"omitempty,url" example:"https://your-s3-bucket.s3.amazonaws.com/course-images/go-programming.jpg"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`

	// Relationships
	Enrollments []Enrollment `json:"enrollments,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Course model
func (Course) TableName() string {
	return "courses"
}

// CourseRequest represents the request payload for creating/updating a course
type CourseRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255" example:"Introduction to Go Programming"`
	Description string  `json:"description" validate:"required,min=1" example:"Learn the fundamentals of Go programming language"`
	Difficulty  string  `json:"difficulty" validate:"required,oneof=Beginner Intermediate Advanced" example:"Beginner"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url" example:"https://your-s3-bucket.s3.amazonaws.com/course-images/go-programming.jpg"`
}

// CourseResponse represents the response payload for course operations
type CourseResponse struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string    `json:"title" example:"Introduction to Go Programming"`
	Description string    `json:"description" example:"Learn the fundamentals of Go programming language"`
	Difficulty  string    `json:"difficulty" example:"Beginner"`
	ImageURL    *string   `json:"image_url,omitempty" example:"https://your-s3-bucket.s3.amazonaws.com/course-images/go-programming.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// ToResponse converts Course model to CourseResponse
func (c *Course) ToResponse() CourseResponse {
	return CourseResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Difficulty:  c.Difficulty,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt,
	}
}
