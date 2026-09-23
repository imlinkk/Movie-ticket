package tests

import (
	"context"
	"testing"
	"time"

	"movie-ticket/internal/models"
	"movie-ticket/internal/services"
	appErrors "movie-ticket/internal/utils/errors"

	"github.com/stretchr/testify/assert"
)

func TestTicketService_BookTicket_ShowtimeNotFound(t *testing.T) {
	mockTicketRepo := new(MockTicketRepository)
	mockShowtimeRepo := new(MockShowtimeRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailService := new(MockEmailService)

	ticketService := services.NewTicketService(mockTicketRepo, mockShowtimeRepo, mockUserRepo, mockEmailService)

	ctx := context.Background()
	req := &models.BookTicketRequest{
		ShowtimeID: 999,
		SeatNumber: "A1",
	}

	mockShowtimeRepo.On("GetByID", ctx, uint(999)).Return(nil, appErrors.ErrNotFound)

	ticket, err := ticketService.BookTicket(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, ticket)
	mockShowtimeRepo.AssertExpectations(t)
}

func TestTicketService_BookTicket_NoAvailableSeats(t *testing.T) {
	mockTicketRepo := new(MockTicketRepository)
	mockShowtimeRepo := new(MockShowtimeRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailService := new(MockEmailService)

	ticketService := services.NewTicketService(mockTicketRepo, mockShowtimeRepo, mockUserRepo, mockEmailService)

	ctx := context.Background()
	req := &models.BookTicketRequest{
		ShowtimeID: 1,
		SeatNumber: "A1",
	}

	showtime := &models.Showtime{
		ID:             1,
		AvailableSeats: 0,
		TotalSeats:     50,
	}

	mockShowtimeRepo.On("GetByID", ctx, uint(1)).Return(showtime, nil)

	ticket, err := ticketService.BookTicket(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, ticket)
	mockShowtimeRepo.AssertExpectations(t)
}

func TestTicketService_BookTicket_SeatAlreadyReserved(t *testing.T) {
	mockTicketRepo := new(MockTicketRepository)
	mockShowtimeRepo := new(MockShowtimeRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailService := new(MockEmailService)

	ticketService := services.NewTicketService(mockTicketRepo, mockShowtimeRepo, mockUserRepo, mockEmailService)

	ctx := context.Background()
	req := &models.BookTicketRequest{
		ShowtimeID: 1,
		SeatNumber: "A1",
	}

	showtime := &models.Showtime{
		ID:             1,
		AvailableSeats: 10,
		TotalSeats:     50,
		Price:          12.5,
	}

	existingTicket := &models.Ticket{
		ID:         5,
		ShowtimeID: 1,
		SeatNumber: "A1",
		Status:     models.TicketStatusBooked,
	}

	mockShowtimeRepo.On("GetByID", ctx, uint(1)).Return(showtime, nil)
	mockTicketRepo.On("GetByShowtimeAndSeat", ctx, uint(1), "A1").Return(existingTicket, nil)

	ticket, err := ticketService.BookTicket(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, ticket)
	mockShowtimeRepo.AssertExpectations(t)
	mockTicketRepo.AssertExpectations(t)
}

func TestTicketService_GetByID_Success(t *testing.T) {
	mockTicketRepo := new(MockTicketRepository)
	mockShowtimeRepo := new(MockShowtimeRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailService := new(MockEmailService)

	ticketService := services.NewTicketService(mockTicketRepo, mockShowtimeRepo, mockUserRepo, mockEmailService)

	ctx := context.Background()
	expectedTicket := &models.Ticket{
		ID:          1,
		UserID:      2,
		ShowtimeID:  3,
		SeatNumber:  "B5",
		Price:       15.0,
		Status:      models.TicketStatusBooked,
		BookingTime: time.Now(),
	}

	mockTicketRepo.On("GetByID", ctx, uint(1)).Return(expectedTicket, nil)

	ticket, err := ticketService.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, uint(1), ticket.ID)
	assert.Equal(t, "B5", ticket.SeatNumber)
	mockTicketRepo.AssertExpectations(t)
}
