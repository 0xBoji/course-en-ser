package router

import (
	"log"
	"sonic-labs/course-enrollment-service/internal/config"
	"sonic-labs/course-enrollment-service/internal/handler"
	"sonic-labs/course-enrollment-service/internal/middleware"
	"sonic-labs/course-enrollment-service/internal/repository"
	"sonic-labs/course-enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Custom logging middleware to ensure logs go to our log file
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Add custom request logging
	r.Use(func(c *gin.Context) {
		log.Printf("API Request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
		log.Printf("API Response: %s %s -> %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	})

	// Initialize repositories
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize Redis service
	redisService := service.NewRedisService(cfg)

	// Test Redis connection
	if err := redisService.Ping(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		redisService = nil // Disable Redis if connection fails
	} else {
		log.Println("Redis connected successfully")
	}

	// Initialize services
	courseService := service.NewCourseService(courseRepo, enrollmentRepo, redisService)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo)
	authService := service.NewAuthService(userRepo)
	studentService := service.NewStudentService(enrollmentRepo)

	// Initialize S3 service
	s3Service := service.NewS3Service()

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService, s3Service)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)
	studentHandler := handler.NewStudentHandler(studentService)
	authHandler := handler.NewAuthHandler(authService)
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status":   "healthy",
			"service":  "course-enrollment-service",
			"database": "connected",
		}

		// Check Redis status
		if redisService != nil {
			if err := redisService.Ping(); err != nil {
				health["redis"] = "disconnected"
				health["redis_error"] = err.Error()
			} else {
				health["redis"] = "connected"
			}
		} else {
			health["redis"] = "disabled"
		}

		c.JSON(200, health)
	})

	// Redis stats endpoint
	r.GET("/redis/stats", func(c *gin.Context) {
		if redisService == nil {
			c.JSON(503, gin.H{
				"error":   "Redis unavailable",
				"message": "Redis service is not available",
			})
			return
		}

		stats, err := redisService.GetStats()
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to get Redis stats",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"data":   stats,
		})
	})

	// API v1 routes - all protected except login
	v1 := r.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)                                                                  // Public - login only
			auth.GET("/profile", middleware.AuthMiddleware(), middleware.AdminMiddleware(), authHandler.GetProfile) // Protected - admin only
		}

		// Public course routes (read-only)
		publicCourses := v1.Group("/courses")
		{
			publicCourses.GET("", courseHandler.GetAllCourses)     // Public - read all courses
			publicCourses.GET("/:id", courseHandler.GetCourseByID) // Public - read specific course
		}

		// Public enrollment routes (read-only)
		publicStudents := v1.Group("/students")
		{
			publicStudents.GET("/:email/enrollments", enrollmentHandler.GetStudentEnrollments) // Public - read student enrollments
		}

		// All other routes require admin authentication
		adminRoutes := v1.Group("")
		adminRoutes.Use(middleware.AdminAuthMiddleware())
		{
			// Course management routes - admin only (write operations)
			courses := adminRoutes.Group("/courses")
			{
				courses.POST("", courseHandler.CreateCourse)                                  // Admin only - create course JSON (default)
				courses.POST("/upload", courseHandler.CreateCourseWithImage)                  // Admin only - create course with image upload
				courses.PUT("/:id", courseHandler.UpdateCourse)                               // Admin only - update course
				courses.DELETE("/:id", courseHandler.DeleteCourse)                            // Admin only - delete course
				courses.GET("/:id/students", courseHandler.GetCourseStudents)                 // Admin only - get course students
				courses.DELETE("/:id/students/:email", courseHandler.RemoveStudentFromCourse) // Admin only - remove student from course
			}

			// Enrollment routes - admin only
			enrollments := adminRoutes.Group("/enrollments")
			{
				enrollments.POST("", enrollmentHandler.EnrollStudent) // Admin only - enroll student
			}

			// Admin routes for student and enrollment management
			admin := adminRoutes.Group("/admin")
			{
				admin.GET("/students", studentHandler.GetAllStudents)             // Admin only - get all students
				admin.GET("/enrollments", studentHandler.GetAllEnrollments)       // Admin only - get all enrollments
				admin.DELETE("/enrollments/:id", studentHandler.DeleteEnrollment) // Admin only - delete enrollment
			}

			// Student management routes - admin only (write operations only, reads are public)
			// Note: Student enrollment reading is available publicly above
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
