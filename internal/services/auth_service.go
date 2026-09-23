package services

import (
	"context"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, appErrors.Conflict("Email is already registered", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.Internal("Failed to hash password", err)
	}

	role := models.RoleUser
	if req.Role == string(models.RoleAdmin) {
		role = models.RoleAdmin
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, appErrors.Internal("Failed to generate token", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.Unauthorized("Invalid email or password", appErrors.ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, appErrors.Unauthorized("Invalid email or password", appErrors.ErrInvalidCredentials)
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, appErrors.Internal("Failed to generate token", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
