package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sonic-labs/course-enrollment-service/internal/config"
	"sonic-labs/course-enrollment-service/internal/constants"
	"sonic-labs/course-enrollment-service/internal/models"

	"github.com/redis/go-redis/v9"
)

// RedisService handles Redis operations
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService creates a new Redis service
func NewRedisService(cfg *config.Config) *RedisService {
	password := cfg.Redis.Password
	// Handle empty password cases
	if password == "none" || password == "empty" || password == "__EMPTY__" || password == "" {
		password = ""
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: password,
		DB:       cfg.Redis.DB,
	})

	return &RedisService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Ping tests Redis connection
func (r *RedisService) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Close closes Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}

// Course caching methods

// SetCourse caches a course
func (r *RedisService) SetCourse(course *models.CourseResponse) error {
	key := fmt.Sprintf("course:%s", course.ID.String())
	data, err := json.Marshal(course)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, constants.CacheTTL).Err()
}

// GetCourse retrieves a cached course
func (r *RedisService) GetCourse(courseID string) (*models.CourseResponse, error) {
	key := fmt.Sprintf("course:%s", courseID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var course models.CourseResponse
	err = json.Unmarshal([]byte(data), &course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

// DeleteCourse removes a course from cache
func (r *RedisService) DeleteCourse(courseID string) error {
	key := fmt.Sprintf("course:%s", courseID)
	return r.client.Del(r.ctx, key).Err()
}

// SetCourses caches all courses list
func (r *RedisService) SetCourses(courses []*models.CourseResponse) error {
	key := "courses:all"
	data, err := json.Marshal(courses)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, constants.CacheTTL).Err()
}

// GetCourses retrieves cached courses list
func (r *RedisService) GetCourses() ([]*models.CourseResponse, error) {
	key := "courses:all"
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var courses []*models.CourseResponse
	err = json.Unmarshal([]byte(data), &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

// InvalidateCoursesCache removes all courses cache
func (r *RedisService) InvalidateCoursesCache() error {
	return r.client.Del(r.ctx, "courses:all").Err()
}

// Session management methods

// SetSession stores a user session
func (r *RedisService) SetSession(sessionID string, userID string, duration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Set(r.ctx, key, userID, duration).Err()
}

// GetSession retrieves a user session
func (r *RedisService) GetSession(sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	userID, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Session not found
		}
		return "", err
	}
	return userID, nil
}

// DeleteSession removes a user session
func (r *RedisService) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Del(r.ctx, key).Err()
}

// Rate limiting methods

// CheckRateLimit checks if a user has exceeded rate limit
func (r *RedisService) CheckRateLimit(userID string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", userID)

	// Get current count
	count, err := r.client.Get(r.ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if count >= limit {
		return false, nil // Rate limit exceeded
	}

	// Increment counter
	pipe := r.client.Pipeline()
	pipe.Incr(r.ctx, key)
	pipe.Expire(r.ctx, key, window)
	_, err = pipe.Exec(r.ctx)

	return err == nil, err
}

// General cache methods

// Set stores a key-value pair with TTL
func (r *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Get retrieves a value by key
func (r *RedisService) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key
func (r *RedisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisService) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// GetStats returns Redis statistics
func (r *RedisService) GetStats() (map[string]interface{}, error) {
	info, err := r.client.Info(r.ctx, "memory", "stats").Result()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"info":      info,
		"db_size":   r.client.DBSize(r.ctx).Val(),
		"ping":      r.client.Ping(r.ctx).Val(),
		"connected": true,
	}

	return stats, nil
}
