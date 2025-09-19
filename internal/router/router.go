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

	// Initialize services
	courseService := service.NewCourseService(courseRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo)
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)
	authHandler := handler.NewAuthHandler(authService)
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "course-enrollment-service",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
		}

		// Course routes
		courses := v1.Group("/courses")
		{
			courses.GET("", courseHandler.GetAllCourses)                                                            // Public - read courses
			courses.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), courseHandler.CreateCourse) // Protected - admin only
			courses.GET("/:id", courseHandler.GetCourseByID)                                                        // Public - read specific course
		}

		// Enrollment routes
		enrollments := v1.Group("/enrollments")
		{
			enrollments.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), enrollmentHandler.EnrollStudent) // Protected - admin only
		}

		// Student routes
		students := v1.Group("/students")
		{
			students.GET("/:email/enrollments", enrollmentHandler.GetStudentEnrollments) // Public - read enrollments
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
