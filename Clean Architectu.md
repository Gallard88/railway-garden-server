# Scaling Your Go API - Clean Architecture Guide

## Problem: Handlers Getting Too Fat

As you add business logic, validation, and complex operations, your handlers become thousands of lines. Here's how to fix it:

## Recommended Architecture (Clean Architecture / Hexagonal)

```
railway-go-api/
├── cmd/
│   └── api/
│       ├── main.go              # App initialization only
│       └── server.go            # HTTP server setup
├── internal/
│   ├── models/                  # Database models (GORM)
│   │   ├── user.go
│   │   ├── post.go
│   │   └── task.go
│   ├── repository/              # Database access layer
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   └── task_repository.go
│   ├── service/                 # Business logic (THE CORE!)
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── auth_service.go
│   │   └── notification_service.go
│   ├── handler/                 # HTTP handlers (thin layer)
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   └── task_handler.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── user_dto.go
│   │   └── post_dto.go
│   ├── validator/               # Custom validation logic
│   │   └── validator.go
│   ├── jobs/                    # Cron jobs
│   │   ├── cleanup_job.go
│   │   └── report_job.go
│   ├── config/                  # Configuration
│   │   └── config.go
│   └── database/               
│       └── db.go
├── pkg/                         # Shared utilities (can be used by other projects)
│   ├── errors/
│   │   └── errors.go
│   ├── logger/
│   │   └── logger.go
│   └── utils/
│       ├── hash.go
│       └── token.go
└── go.mod
```

## The Layers Explained

### 1. **Models** (`internal/models/`)
- Just structs with GORM tags
- No business logic
- Database schema representation

### 2. **Repository** (`internal/repository/`)
- Database queries ONLY
- No business logic
- Returns models or errors
- Can be mocked for testing

### 3. **Service** (`internal/service/`)
- **THIS IS WHERE BUSINESS LOGIC LIVES** ⭐
- Orchestrates repositories
- Validates business rules
- Handles complex operations
- Transactions across multiple repos

### 4. **Handlers** (`internal/handler/`)
- Parse HTTP requests
- Call services
- Format HTTP responses
- Keep them thin (50-100 lines max)

### 5. **DTOs** (`internal/dto/`)
- Request/Response structures
- Different from models (security!)
- Validation tags

---

## Example Implementation

### internal/models/user.go
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `gorm:"not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"` // Never expose in JSON
	Active    bool           `gorm:"default:true" json:"active"`
	Role      string         `gorm:"default:'user'" json:"role"` // user, admin, etc.
}
```

### internal/dto/user_dto.go
```go
package dto

// CreateUserRequest - What we accept from API
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserRequest - Partial updates
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Name   string `json:"name" binding:"omitempty,min=2,max=100"`
	Active *bool  `json:"active"`
}

// UserResponse - What we return (no password!)
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ListUsersResponse - Paginated list
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}
```

### internal/repository/user_repository.go
```go
package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"goapi.railway.app/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	CountActive(ctx context.Context) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, totalCount, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("active = ?", true).
		Count(&count).Error
	return count, err
}
```

### internal/service/user_service.go
```go
package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id uint) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error)
	UpdateUser(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	ActivateUser(ctx context.Context, id uint) error
	DeactivateUser(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
	// Could inject other services here
	// emailService EmailService
	// auditService AuditService
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser - Business logic for creating users
func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 1. Business Rule: Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// 2. Business Rule: Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create model
	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		Active:   true,
		Role:     "user",
	}

	// 4. Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. Business Logic: Could send welcome email here
	// s.emailService.SendWelcome(user.Email)

	// 6. Business Logic: Could log audit event
	// s.auditService.LogUserCreated(user.ID)

	// 7. Convert to DTO and return
	return s.modelToResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.modelToResponse(user), nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error) {
	// Business Rule: Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, totalCount, err := s.userRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.modelToResponse(&user)
	}

	return &dto.ListUsersResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 1. Get existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Business Rule: If changing email, check it's not taken
	if req.Email != "" && req.Email != user.Email {
		existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
		if existing != nil {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	// 3. Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	// 4. Save
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.modelToResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	// Business Rule: Could check if user has dependencies
	// Business Rule: Could anonymize data instead of deleting
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ActivateUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Active {
		return errors.New("user already active")
	}

	user.Active = true
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeactivateUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !user.Active {
		return errors.New("user already inactive")
	}

	user.Active = false
	return s.userRepo.Update(ctx, user)
}

// Helper to convert model to DTO
func (s *userService) modelToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Active:    user.Active,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
```

### internal/handler/user_handler.go
```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser - THIN handler, just HTTP stuff
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	// 1. Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call service (where business logic lives)
	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		// Could have better error handling here
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// RegisterRoutes - Register all user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", h.ListUsers)
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}
```

### cmd/api/main.go (Updated)
```go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"goapi.railway.app/internal/database"
	"goapi.railway.app/internal/handler"
	"goapi.railway.app/internal/repository"
	"goapi.railway.app/internal/service"
)

const version = "0.0.2"

func main() {
	godotenv.Load()

	// Get port
	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 4000
	}

	// Connect to database
	db, err := database.New()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	// postRepo := repository.NewPostRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	// postService := service.NewPostService(postRepo, userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	// postHandler := handler.NewPostHandler(postService)

	// Initialize cron
	cronScheduler := cron.New()
	// Setup cron jobs...
	cronScheduler.Start()
	defer cronScheduler.Stop()

	// Setup Gin
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()

	// Register routes
	v1 := router.Group("/v1")
	{
		router.GET("/v1/healthcheck", healthcheck)
		userHandler.RegisterRoutes(v1)
		// postHandler.RegisterRoutes(v1)
	}

	// Start server
	addr := fmt.Sprintf(":%d", intPort)
	log.Printf("Server starting on %s (version %s)", addr, version)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func healthcheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "version": version})
}
```

---

## When to Use Each Layer

### Repository Layer
```go
// ✅ GOOD - Just data access
func (r *userRepository) FindActiveUsers() ([]User, error)

// ❌ BAD - Business logic in repository
func (r *userRepository) FindUsersAndSendEmail() error
```

### Service Layer
```go
// ✅ GOOD - Business logic
func (s *userService) RegisterUser(req CreateUserRequest) error {
    // Validate email not taken
    // Hash password
    // Create user
    // Send welcome email
    // Log audit event
}

// ❌ BAD - Just wrapping repository
func (s *userService) GetUser(id uint) (*User, error) {
    return s.repo.FindByID(id) // Too simple, just call repo from handler
}
```

### Handler Layer
```go
// ✅ GOOD - Thin HTTP layer
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.BindJSON(&req)
    user, err := h.service.CreateUser(req)
    c.JSON(200, user)
}

// ❌ BAD - Business logic in handler
func (h *Handler) CreateUser(c *gin.Context) {
    // Validation
    // Password hashing
    // Database save
    // Email sending
    // 100+ lines of code
}
```

---

## Complex Example: Post Creation with Multiple Services

```go
// service/post_service.go
type PostService interface {
	CreatePost(ctx context.Context, req dto.CreatePostRequest, userID uint) (*dto.PostResponse, error)
}

type postService struct {
	postRepo         repository.PostRepository
	userRepo         repository.UserRepository
	notificationSvc  NotificationService // Another service!
	searchIndexer    SearchIndexer       // External service
}

func (s *postService) CreatePost(ctx context.Context, req dto.CreatePostRequest, userID uint) (*dto.PostResponse, error) {
	// 1. Business Rule: Verify user exists and is active
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.Active {
		return nil, errors.New("user account is inactive")
	}

	// 2. Business Rule: Check content moderation
	if containsProfanity(req.Content) {
		return nil, errors.New("content violates community guidelines")
	}

	// 3. Business Rule: Premium users get more features
	if req.Featured && user.Role != "premium" {
		return nil, errors.New("featured posts are premium-only")
	}

	// 4. Create post
	post := &models.Post{
		Title:     req.Title,
		Content:   req.Content,
		UserID:    userID,
		Published: req.Published,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	// 5. Business Logic: Index for search (async would be better)
	go s.searchIndexer.IndexPost(post)

	// 6. Business Logic: Notify followers
	if req.Published {
		go s.notificationSvc.NotifyFollowers(userID, post.ID)
	}

	return s.modelToResponse(post), nil
}
```

---

## Testing Benefits

With this structure, testing becomes easy:

```go
// service/user_service_test.go
func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	// Mock repository
	mockRepo := &MockUserRepository{
		FindByEmailFunc: func(email string) (*models.User, error) {
			return &models.User{Email: email}, nil // Email exists!
		},
	}

	service := NewUserService(mockRepo)

	req := dto.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test",
	}

	// Should fail because email exists
	_, err := service.CreateUser(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already in use")
}
```

---

## Summary: Where Does Everything Go?

| Concern | Layer | Example |
|---------|-------|---------|
| HTTP parsing | Handler | `c.BindJSON()`, `c.Param("id")` |
| Validation | DTO + Service | Struct tags + business rules |
| Business rules | Service | "Email must be unique", "Premium only" |
| Database queries | Repository | `db.Where().Find()` |
| Transactions | Service | Coordinate multiple repos |
| External APIs | Service | Call other services |
| Password hashing | Service | Before saving user |
| Email sending | Service (or dedicated service) | After user creation |
| Cron jobs | Jobs (calling Services) | Use services, not repos directly |

**Golden Rule:** If it's not HTTP-specific, it doesn't belong in the handler!

This keeps your code:
- ✅ Testable
- ✅ Reusable (same service for HTTP + CLI + gRPC)
- ✅ Maintainable (easy to find things)
- ✅ Scalable (add features without spaghetti)