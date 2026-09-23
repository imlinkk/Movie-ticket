package tests

import (
	"context"
	"testing"

	"movie-ticket/config"
	"movie-ticket/internal/models"
	"movie-ticket/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	config.AppConfig = &config.Config{
		AppEnv:     "test",
		JWTSecret:  "test-secret-key-1234567890123456",
		JWTExpireH: 24,
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     "user",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	res, err := authService.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, req.Email, res.User.Email)
	assert.Equal(t, req.Name, res.User.Name)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "Duplicate User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &models.User{
		ID:    1,
		Name:  "Existing",
		Email: "existing@example.com",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	res, err := authService.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()
	password := "correctPassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:       10,
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: string(hashedPassword),
		Role:     models.RoleUser,
	}

	mockRepo.On("GetByEmail", ctx, "alice@example.com").Return(user, nil)

	req := &models.LoginRequest{
		Email:    "alice@example.com",
		Password: password,
	}

	res, err := authService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, uint(10), res.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rightPassword"), bcrypt.DefaultCost)

	user := &models.User{
		ID:       10,
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByEmail", ctx, "alice@example.com").Return(user, nil)

	req := &models.LoginRequest{
		Email:    "alice@example.com",
		Password: "wrongPassword",
	}

	res, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
