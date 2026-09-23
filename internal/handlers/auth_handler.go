package handlers

import (
	"net/http"

	"movie-ticket/internal/middleware"
	"movie-ticket/internal/models"
	"movie-ticket/internal/services"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new account with email, name, password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User Registration Info"
// @Success 201 {object} models.StandardResponse{data=models.AuthResponse}
// @Failure 400 {object} models.StandardResponse
// @Failure 409 {object} models.StandardResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Credentials"
// @Success 200 {object} models.StandardResponse{data=models.AuthResponse}
// @Failure 400 {object} models.StandardResponse
// @Failure 401 {object} models.StandardResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get authenticated user profile from token
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.StandardResponse{data=models.User}
// @Failure 401 {object} models.StandardResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	val, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Error(c, appErrors.Unauthorized("Unauthorized", nil))
		return
	}

	userID := val.(uint)
	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved successfully", user)
}
