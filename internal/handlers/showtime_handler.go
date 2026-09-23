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

type ShowtimeHandler struct {
	showtimeService services.ShowtimeService
}

func NewShowtimeHandler(showtimeService services.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{showtimeService: showtimeService}
}

// ListShowtimes godoc
// @Summary List showtimes
// @Description Get paginated list of showtimes with filters for movie, cinema, or date (YYYY-MM-DD)
// @Tags Showtimes
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param movie_id query int false "Filter by Movie ID"
// @Param cinema_id query int false "Filter by Cinema ID"
// @Param date query string false "Filter by date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Showtime}
// @Router /showtimes [get]
func (h *ShowtimeHandler) List(c *gin.Context) {
	params := pagination.GetPaginationParams(c)
	movieID, _ := strconv.ParseUint(c.Query("movie_id"), 10, 32)
	cinemaID, _ := strconv.ParseUint(c.Query("cinema_id"), 10, 32)
	dateStr := c.Query("date")

	showtimes, total, err := h.showtimeService.List(c.Request.Context(), params, uint(movieID), uint(cinemaID), dateStr)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "Showtimes retrieved successfully", showtimes, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetByID godoc
// @Summary Get showtime by ID
// @Description Retrieve showtime details including available seats
// @Tags Showtimes
// @Produce json
// @Param id path int true "Showtime ID"
// @Success 200 {object} models.StandardResponse{data=models.Showtime}
// @Failure 404 {object} models.StandardResponse
// @Router /showtimes/{id} [get]
func (h *ShowtimeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid showtime ID", err))
		return
	}

	showtime, err := h.showtimeService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Showtime retrieved successfully", showtime)
}

// Create godoc
// @Summary Create showtime
// @Description Create a showtime for a movie and cinema (Admin only)
// @Tags Showtimes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateShowtimeRequest true "Showtime Info"
// @Success 201 {object} models.StandardResponse{data=models.Showtime}
// @Failure 400 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /showtimes [post]
func (h *ShowtimeHandler) Create(c *gin.Context) {
	var req models.CreateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	showtime, err := h.showtimeService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Showtime created successfully", showtime)
}

// Update godoc
// @Summary Update showtime
// @Description Update showtime parameters (Admin only)
// @Tags Showtimes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Showtime ID"
// @Param request body models.UpdateShowtimeRequest true "Update Info"
// @Success 200 {object} models.StandardResponse{data=models.Showtime}
// @Failure 400 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /showtimes/{id} [put]
func (h *ShowtimeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid showtime ID", err))
		return
	}

	var req models.UpdateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	showtime, err := h.showtimeService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Showtime updated successfully", showtime)
}

// Delete godoc
// @Summary Delete showtime
// @Description Soft delete a showtime (Admin only)
// @Tags Showtimes
// @Security BearerAuth
// @Produce json
// @Param id path int true "Showtime ID"
// @Success 200 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /showtimes/{id} [delete]
func (h *ShowtimeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid showtime ID", err))
		return
	}

	if err := h.showtimeService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Showtime deleted successfully", nil)
}
