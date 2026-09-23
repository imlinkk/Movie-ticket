package handlers

import (
	"net/http"
	"strconv"

	"movie-ticket/internal/models"
	"movie-ticket/internal/services"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"
	"movie-ticket/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type CinemaHandler struct {
	cinemaService services.CinemaService
}

func NewCinemaHandler(cinemaService services.CinemaService) *CinemaHandler {
	return &CinemaHandler{cinemaService: cinemaService}
}

// ListCinemas godoc
// @Summary List cinemas
// @Description Get paginated list of cinemas with optional city filter
// @Tags Cinemas
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param city query string false "Filter by city"
// @Param search query string false "Search by name or address"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Cinema}
// @Router /cinemas [get]
func (h *CinemaHandler) List(c *gin.Context) {
	params := pagination.GetPaginationParams(c)
	city := c.Query("city")

	cinemas, total, err := h.cinemaService.List(c.Request.Context(), params, city)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "Cinemas retrieved successfully", cinemas, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetByID godoc
// @Summary Get cinema by ID
// @Description Retrieve cinema details
// @Tags Cinemas
// @Produce json
// @Param id path int true "Cinema ID"
// @Success 200 {object} models.StandardResponse{data=models.Cinema}
// @Failure 404 {object} models.StandardResponse
// @Router /cinemas/{id} [get]
func (h *CinemaHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid cinema ID", err))
		return
	}

	cinema, err := h.cinemaService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Cinema retrieved successfully", cinema)
}

// Create godoc
// @Summary Create cinema
// @Description Create a cinema (Admin only)
// @Tags Cinemas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateCinemaRequest true "Cinema data"
// @Success 201 {object} models.StandardResponse{data=models.Cinema}
// @Failure 400 {object} models.StandardResponse
// @Failure 401 {object} models.StandardResponse
// @Failure 403 {object} models.StandardResponse
// @Router /cinemas [post]
func (h *CinemaHandler) Create(c *gin.Context) {
	var req models.CreateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	cinema, err := h.cinemaService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Cinema created successfully", cinema)
}

// Update godoc
// @Summary Update cinema
// @Description Update cinema details (Admin only)
// @Tags Cinemas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cinema ID"
// @Param request body models.UpdateCinemaRequest true "Update Cinema payload"
// @Success 200 {object} models.StandardResponse{data=models.Cinema}
// @Failure 400 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /cinemas/{id} [put]
func (h *CinemaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid cinema ID", err))
		return
	}

	var req models.UpdateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	cinema, err := h.cinemaService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Cinema updated successfully", cinema)
}

// Delete godoc
// @Summary Delete cinema
// @Description Soft delete cinema (Admin only)
// @Tags Cinemas
// @Security BearerAuth
// @Produce json
// @Param id path int true "Cinema ID"
// @Success 200 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /cinemas/{id} [delete]
func (h *CinemaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid cinema ID", err))
		return
	}

	if err := h.cinemaService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Cinema deleted successfully", nil)
}
