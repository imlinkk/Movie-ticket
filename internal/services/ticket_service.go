package services

import (
	"context"
	"time"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"
)

type TicketService interface {
	BookTicket(ctx context.Context, userID uint, req *models.BookTicketRequest) (*models.Ticket, error)
	GetByID(ctx context.Context, id uint) (*models.Ticket, error)
	GetMyTickets(ctx context.Context, userID uint, params pagination.Params) ([]models.Ticket, int64, error)
	GetAllTickets(ctx context.Context, params pagination.Params, status string) ([]models.Ticket, int64, error)
	CancelTicket(ctx context.Context, id uint, userID uint, isAdmin bool) (*models.Ticket, error)
	UpdateStatus(ctx context.Context, id uint, status models.TicketStatus) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo   repositories.TicketRepository
	showtimeRepo repositories.ShowtimeRepository
	userRepo     repositories.UserRepository
	emailService EmailService
}

func NewTicketService(
	ticketRepo repositories.TicketRepository,
	showtimeRepo repositories.ShowtimeRepository,
	userRepo repositories.UserRepository,
	emailService EmailService,
) TicketService {
	return &ticketService{
		ticketRepo:   ticketRepo,
		showtimeRepo: showtimeRepo,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *ticketService) BookTicket(ctx context.Context, userID uint, req *models.BookTicketRequest) (*models.Ticket, error) {
	// 1. Verify showtime exists
	showtime, err := s.showtimeRepo.GetByID(ctx, req.ShowtimeID)
	if err != nil {
		return nil, appErrors.NotFound("Showtime not found", err)
	}

	if showtime.AvailableSeats <= 0 {
		return nil, appErrors.Conflict("Showtime has no available seats remaining", appErrors.ErrNoAvailableSeats)
	}

	// 2. Check if seat is already booked
	existingTicket, err := s.ticketRepo.GetByShowtimeAndSeat(ctx, req.ShowtimeID, req.SeatNumber)
	if err != nil {
		return nil, err
	}
	if existingTicket != nil {
		return nil, appErrors.Conflict("Seat "+req.SeatNumber+" is already reserved", appErrors.ErrSeatAlreadyBooked)
	}

	// 3. Begin DB Transaction for atomic reservation
	db := s.ticketRepo.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, appErrors.Internal("Failed to begin transaction", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Decrement available seat with check
	if err := s.showtimeRepo.DecrementAvailableSeats(ctx, tx, req.ShowtimeID); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create ticket record
	ticket := &models.Ticket{
		UserID:      userID,
		ShowtimeID:  req.ShowtimeID,
		SeatNumber:  req.SeatNumber,
		Price:       showtime.Price,
		Status:      models.TicketStatusBooked,
		BookingTime: time.Now(),
	}

	if err := s.ticketRepo.CreateWithTx(ctx, tx, ticket); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, appErrors.Internal("Failed to commit booking transaction", err)
	}

	// Fetch complete populated ticket
	finalTicket, err := s.ticketRepo.GetByID(ctx, ticket.ID)
	if err != nil {
		return ticket, nil
	}

	// 4. Asynchronously trigger background confirmation email
	user, err := s.userRepo.GetByID(ctx, userID)
	if err == nil && s.emailService != nil {
		s.emailService.SendTicketConfirmation(user, finalTicket, showtime)
	}

	return finalTicket, nil
}

func (s *ticketService) GetByID(ctx context.Context, id uint) (*models.Ticket, error) {
	return s.ticketRepo.GetByID(ctx, id)
}

func (s *ticketService) GetMyTickets(ctx context.Context, userID uint, params pagination.Params) ([]models.Ticket, int64, error) {
	return s.ticketRepo.ListByUser(ctx, userID, params)
}

func (s *ticketService) GetAllTickets(ctx context.Context, params pagination.Params, status string) ([]models.Ticket, int64, error) {
	return s.ticketRepo.ListAll(ctx, params, status)
}

func (s *ticketService) CancelTicket(ctx context.Context, id uint, userID uint, isAdmin bool) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !isAdmin && ticket.UserID != userID {
		return nil, appErrors.Forbidden("You are not allowed to cancel this ticket", nil)
	}

	if ticket.Status == models.TicketStatusCancelled {
		return nil, appErrors.BadRequest("Ticket is already cancelled", nil)
	}

	tx := s.ticketRepo.GetDB().Begin()
	if tx.Error != nil {
		return nil, appErrors.Internal("Failed to begin transaction", tx.Error)
	}

	ticket.Status = models.TicketStatusCancelled
	if err := tx.Save(ticket).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Restore seat availability
	if err := s.showtimeRepo.IncrementAvailableSeats(ctx, tx, ticket.ShowtimeID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, appErrors.Internal("Failed to commit ticket cancellation", err)
	}

	return ticket, nil
}

func (s *ticketService) UpdateStatus(ctx context.Context, id uint, status models.TicketStatus) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ticket.Status = status
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
