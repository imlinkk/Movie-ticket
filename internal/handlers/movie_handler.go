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

type MovieHandler struct {
	movieService services.MovieService
}

func NewMovieHandler(movieService services.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

// ListMovies godoc
// @Summary List movies
// @Description Get paginated list of movies with optional search, sorting, and genre filter
// @Tags Movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param genre query string false "Filter by genre"
// @Param search query string false "Search by title or description"
// @Param sort_by query string false "Field to sort by (e.g. release_date, rating, title)" default(created_at)
// @Param order query string false "Order direction (asc/desc)" default(desc)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Movie}
// @Failure 500 {object} models.StandardResponse
// @Router /movies [get]
func (h *MovieHandler) List(c *gin.Context) {
	params := pagination.GetPaginationParams(c)
	genre := c.Query("genre")

	movies, total, err := h.movieService.List(c.Request.Context(), params, genre)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "Movies retrieved successfully", movies, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetByID godoc
// @Summary Get movie by ID
// @Description Retrieve movie details by ID
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.StandardResponse{data=models.Movie}
// @Failure 404 {object} models.StandardResponse
// @Router /movies/{id} [get]
func (h *MovieHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid movie ID", err))
		return
	}

	movie, err := h.movieService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Movie retrieved successfully", movie)
}

// Create godoc
// @Summary Create a new movie
// @Description Create a movie entry (Admin only)
// @Tags Movies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateMovieRequest true "Movie details"
// @Success 201 {object} models.StandardResponse{data=models.Movie}
// @Failure 400 {object} models.StandardResponse
// @Failure 401 {object} models.StandardResponse
// @Failure 403 {object} models.StandardResponse
// @Router /movies [post]
func (h *MovieHandler) Create(c *gin.Context) {
	var req models.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	movie, err := h.movieService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Movie created successfully", movie)
}

// Update godoc
// @Summary Update movie
// @Description Update movie details by ID (Admin only)
// @Tags Movies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param request body models.UpdateMovieRequest true "Update payload"
// @Success 200 {object} models.StandardResponse{data=models.Movie}
// @Failure 400 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /movies/{id} [put]
func (h *MovieHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid movie ID", err))
		return
	}

	var req models.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	movie, err := h.movieService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Movie updated successfully", movie)
}

// Delete godoc
// @Summary Delete movie
// @Description Soft delete a movie by ID (Admin only)
// @Tags Movies
// @Security BearerAuth
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.StandardResponse
// @Failure 404 {object} models.StandardResponse
// @Router /movies/{id} [delete]
func (h *MovieHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid movie ID", err))
		return
	}

	if err := h.movieService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Movie deleted successfully", nil)
}
