package handlers

import (
	"net/http"
	"strconv"

	"movie-ticket/internal/middleware"
	"movie-ticket/internal/models"
	"movie-ticket/internal/services"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"
	"movie-ticket/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users (Admin only)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param search query string false "Search name or email"
// @Success 200 {object} models.PaginatedResponse{data=[]models.User}
// @Failure 401 {object} models.StandardResponse
// @Failure 403 {object} models.StandardResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	params := pagination.GetPaginationParams(c)
	users, total, err := h.userService.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "Users retrieved successfully", users, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.StandardResponse{data=models.User}
// @Failure 404 {object} models.StandardResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid user ID", err))
		return
	}

	currentUserID := c.GetUint(middleware.CtxUserIDKey)
	currentUserRole := c.GetString(middleware.CtxUserRole)

	// User can only view self unless admin
	if currentUserRole != string(models.RoleAdmin) && currentUserID != uint(id) {
		response.Error(c, appErrors.Forbidden("Access denied to view this user", nil))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// Update godoc
// @Summary Update user
// @Description Update user details
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "Update Data"
// @Success 200 {object} models.StandardResponse{data=models.User}
// @Failure 400 {object} models.StandardResponse
// @Failure 403 {object} models.StandardResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid user ID", err))
		return
	}

	currentUserID := c.GetUint(middleware.CtxUserIDKey)
	currentUserRole := c.GetString(middleware.CtxUserRole)

	// Only admin can change roles or update other users
	if currentUserRole != string(models.RoleAdmin) && currentUserID != uint(id) {
		response.Error(c, appErrors.Forbidden("Access denied to update this user", nil))
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Normal user cannot elevate role
	if currentUserRole != string(models.RoleAdmin) {
		req.Role = ""
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// Delete godoc
// @Summary Delete user
// @Description Soft delete user by ID (Admin only)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid user ID", err))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
