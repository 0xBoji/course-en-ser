# Course Enrollment Service

A robust backend API service for managing course catalog and student enrollments for the Sonic University platform. Built with Go, Gin framework, PostgreSQL, and comprehensive testing.

## Features

- **RESTful API** with four core endpoints for course and enrollment management
- **PostgreSQL Database** with GORM ORM for robust data persistence
- **Swagger Documentation** with interactive API explorer
- **Comprehensive Testing** including unit and integration tests
- **Docker Support** for easy deployment and development
- **Input Validation** with proper error handling and HTTP status codes
- **Duplicate Prevention** for student enrollments with database constraints
- **Email Validation** for student enrollment requests

## API Endpoints

### Courses
- `GET /api/v1/courses` - Retrieve all available courses
- `POST /api/v1/courses` - Create a new course

### Enrollments
- `POST /api/v1/enrollments` - Enroll a student in a course
- `GET /api/v1/students/:email/enrollments` - Get all courses a student is enrolled in

### Documentation
- `GET /swagger/index.html` - Interactive Swagger API documentation
- `GET /health` - Health check endpoint

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (optional)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd course-enrollment-service
   ```

2. **Set up configuration**
   ```bash
   # Copy environment file
   cp .env.example .env

   # Edit .env with your actual RDS database credentials
   # IMPORTANT: Never commit .env to git!
   ```

   **Configuration Method:**
   - Environment variables (`.env` file for development)
   - Environment variables (GitHub Secrets for CI/CD)
   - Default values (fallback)

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb course_enrollment
   
   # Run migrations (optional - app will auto-migrate)
   psql -d course_enrollment -f migrations/001_create_courses_table.sql
   psql -d course_enrollment -f migrations/002_create_enrollments_table.sql
   psql -d course_enrollment -f migrations/003_seed_demo_courses.sql
   ```

5. **Generate Swagger documentation**
   ```bash
   make swagger
   ```

6. **Run the application**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080` and Swagger documentation at `http://localhost:8080/swagger/index.html`.

### Docker Setup

1. **Run with Docker Compose**
   ```bash
   make docker-up
   ```

This will start both the application and PostgreSQL database in containers.

## Configuration

The application can be configured using environment variables or a `config/app.ini` file:

### Environment Variables

- `PORT` - Server port (default: 8080)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: course_enrollment)
- `DB_SSLMODE` - SSL mode (default: disable)

## API Usage Examples

### Create a Course

```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Go Programming",
    "description": "Learn the fundamentals of Go programming language",
    "difficulty": "Beginner"
  }'
```

### Get All Courses

```bash
curl http://localhost:8080/api/v1/courses
```

### Enroll a Student

```bash
curl -X POST http://localhost:8080/api/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_email": "student@example.com",
    "course_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

### Get Student Enrollments

```bash
curl http://localhost:8080/api/v1/students/student@example.com/enrollments
```

## Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

### Integration Tests

Integration tests require a PostgreSQL database. Set up a test database:

```bash
createdb course_enrollment_test
```

The integration tests will automatically handle schema setup and cleanup.

## Database Schema

### Courses Table

- `id` (UUID, Primary Key)
- `title` (VARCHAR, NOT NULL)
- `description` (TEXT, NOT NULL)
- `difficulty` (VARCHAR, CHECK: Beginner/Intermediate/Advanced)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Enrollments Table

- `id` (UUID, Primary Key)
- `student_email` (VARCHAR, NOT NULL)
- `course_id` (UUID, Foreign Key to courses.id)
- `enrolled_at` (TIMESTAMP)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- **Unique Constraint**: (student_email, course_id) to prevent duplicates

## Development

### Available Make Commands

- `make build` - Build the application
- `make run` - Build and run the application
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make swagger` - Generate Swagger documentation
- `make docker-build` - Build Docker image
- `make docker-up` - Run with Docker Compose
- `make docker-down` - Stop Docker Compose
- `make fmt` - Format code
- `make clean` - Clean build artifacts

### Project Structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and migrations
│   ├── handler/        # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   ├── router/         # HTTP routing
│   └── service/        # Business logic layer
├── tests/              # Integration tests
├── migrations/         # SQL migration scripts
├── docs/               # Generated Swagger documentation
└── Makefile           # Build automation
```

## Error Handling

The API returns standardized error responses:

```json
{
  "error": "Validation failed",
  "message": "Title is required"
}
```

### HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `409` - Conflict (duplicate enrollment)
- `500` - Internal Server Error

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.
# CI/CD Test Fri Sep 19 17:49:43 +07 2025
