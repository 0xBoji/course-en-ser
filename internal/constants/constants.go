package constants

import "time"

// HTTP Status Messages
const (
	// Success Messages
	MsgSuccess = "Success"

	// Error Messages
	MsgInternalServerError = "Internal server error"
	MsgBadRequest          = "Bad request"
	MsgUnauthorized        = "Unauthorized"
	MsgForbidden           = "Forbidden"
	MsgNotFound            = "Not found"
	MsgConflict            = "Conflict"

	// Authentication Messages
	MsgInvalidCredentials  = "Invalid username or password"
	MsgAuthHeaderRequired  = "Authorization header is required"
	MsgInvalidTokenFormat  = "Invalid token format"
	MsgJWTTokenInvalid     = "JWT token is invalid or expired"
	MsgAdminAccessRequired = "Admin access required"

	// Course Messages
	MsgCourseNotFound        = "The requested course does not exist"
	MsgCourseCreated         = "Course created successfully"
	MsgInvalidCourseIDFormat = "Invalid course ID format"

	// Enrollment Messages
	MsgEnrollmentCreated      = "Student enrolled successfully"
	MsgStudentAlreadyEnrolled = "Student is already enrolled in this course"
	MsgCourseNotExist         = "The specified course does not exist"
	MsgInvalidEmailFormat     = "Invalid email format"

	// Validation Messages
	MsgTitleRequired       = "Title is required"
	MsgDescriptionRequired = "Description is required"
	MsgDifficultyRequired  = "Difficulty is required"
	MsgDifficultyInvalid   = "Difficulty must be one of: Beginner, Intermediate, Advanced"
	MsgEmailRequired       = "Student email is required"
	MsgCourseIDRequired    = "Course ID is required"
)

// JWT Constants
const (
	JWTTokenExpiry = 24 * time.Hour
	JWTIssuer      = "sonic-labs-course-enrollment"
)

// Cache Constants
const (
	CacheTTL          = 15 * time.Minute
	SessionTTL        = 24 * time.Hour
	RateLimitWindow   = 1 * time.Minute
	RateLimitRequests = 60
)

// User Roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Database Table Names
const (
	TableUsers       = "users"
	TableCourses     = "courses"
	TableEnrollments = "enrollments"
)

// Course Difficulty Levels
const (
	DifficultyBeginner     = "Beginner"
	DifficultyIntermediate = "Intermediate"
	DifficultyAdvanced     = "Advanced"
)

// Pagination settings
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	MinPageSize     = 1
)

// HTTP Headers
const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
)

// Content Types
const (
	ContentTypeJSON = "application/json"
)

// API Paths
const (
	APIBasePath = "/api/v1"

	// Auth paths
	AuthBasePath = APIBasePath + "/auth"
	AuthLogin    = "/login"
	AuthProfile  = "/profile"

	// Course paths
	CoursesBasePath = APIBasePath + "/courses"
	CourseByID      = "/:id"

	// Enrollment paths
	EnrollmentsBasePath = APIBasePath + "/enrollments"
	StudentEnrollments  = APIBasePath + "/students/:email/enrollments"
)
