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

type TicketHandler struct {
	ticketService services.TicketService
}

func NewTicketHandler(ticketService services.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

// BookTicket godoc
// @Summary Book a movie ticket
// @Description Book a seat for a specific showtime. Uses database transaction to prevent double booking.
// @Tags Tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.BookTicketRequest true "Booking Info"
// @Success 201 {object} models.StandardResponse{data=models.Ticket}
// @Failure 400 {object} models.StandardResponse
// @Failure 409 {object} models.StandardResponse
// @Router /tickets/book [post]
func (h *TicketHandler) BookTicket(c *gin.Context) {
	userID := c.GetUint(middleware.CtxUserIDKey)

	var req models.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ticket, err := h.ticketService.BookTicket(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Ticket booked successfully! A confirmation email is on its way.", ticket)
}

// GetMyTickets godoc
// @Summary Get current user's tickets
// @Description Get paginated list of tickets belonging to authenticated user
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Ticket}
// @Router /tickets/my-tickets [get]
func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID := c.GetUint(middleware.CtxUserIDKey)
	params := pagination.GetPaginationParams(c)

	tickets, total, err := h.ticketService.GetMyTickets(c.Request.Context(), userID, params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "My tickets retrieved successfully", tickets, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetAllTickets godoc
// @Summary List all tickets
// @Description Get paginated list of all tickets across the system (Admin only)
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param status query string false "Filter by status (booked, paid, cancelled)"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Ticket}
// @Router /tickets [get]
func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	params := pagination.GetPaginationParams(c)
	status := c.Query("status")

	tickets, total, err := h.ticketService.GetAllTickets(c.Request.Context(), params, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, "Tickets retrieved successfully", tickets, pagination.BuildPagination(params.Page, params.Limit, total))
}

// GetByID godoc
// @Summary Get ticket by ID
// @Description Retrieve single ticket details
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.StandardResponse{data=models.Ticket}
// @Failure 404 {object} models.StandardResponse
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid ticket ID", err))
		return
	}

	ticket, err := h.ticketService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	currentUserID := c.GetUint(middleware.CtxUserIDKey)
	role := c.GetString(middleware.CtxUserRole)
	if role != string(models.RoleAdmin) && ticket.UserID != currentUserID {
		response.Error(c, appErrors.Forbidden("Access denied to view this ticket", nil))
		return
	}

	response.Success(c, http.StatusOK, "Ticket retrieved successfully", ticket)
}

// CancelTicket godoc
// @Summary Cancel ticket
// @Description Cancel a booked ticket and restore showtime seat availability
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.StandardResponse{data=models.Ticket}
// @Failure 400 {object} models.StandardResponse
// @Failure 403 {object} models.StandardResponse
// @Router /tickets/{id}/cancel [put]
func (h *TicketHandler) CancelTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid ticket ID", err))
		return
	}

	currentUserID := c.GetUint(middleware.CtxUserIDKey)
	role := c.GetString(middleware.CtxUserRole)
	isAdmin := role == string(models.RoleAdmin)

	ticket, err := h.ticketService.CancelTicket(c.Request.Context(), uint(id), currentUserID, isAdmin)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Ticket cancelled successfully", ticket)
}

// UpdateStatus godoc
// @Summary Update ticket status
// @Description Update status of a ticket (Admin only)
// @Tags Tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param request body models.UpdateTicketStatusRequest true "Status update"
// @Success 200 {object} models.StandardResponse{data=models.Ticket}
// @Router /tickets/{id}/status [patch]
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, appErrors.BadRequest("Invalid ticket ID", err))
		return
	}

	var req models.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ticket, err := h.ticketService.UpdateStatus(c.Request.Context(), uint(id), models.TicketStatus(req.Status))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Ticket status updated successfully", ticket)
}
