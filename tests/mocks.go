package tests

import (
	"context"

	"movie-ticket/internal/models"
	"movie-ticket/internal/utils/pagination"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, params pagination.Params) ([]models.User, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

// MockMovieRepository
type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MockMovieRepository) GetByID(ctx context.Context, id uint) (*models.Movie, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), args.Error(1)
}

func (m *MockMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MockMovieRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMovieRepository) List(ctx context.Context, params pagination.Params, genre string) ([]models.Movie, int64, error) {
	args := m.Called(ctx, params, genre)
	return args.Get(0).([]models.Movie), args.Get(1).(int64), args.Error(2)
}

// MockShowtimeRepository
type MockShowtimeRepository struct {
	mock.Mock
}

func (m *MockShowtimeRepository) Create(ctx context.Context, showtime *models.Showtime) error {
	args := m.Called(ctx, showtime)
	return args.Error(0)
}

func (m *MockShowtimeRepository) GetByID(ctx context.Context, id uint) (*models.Showtime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) Update(ctx context.Context, showtime *models.Showtime) error {
	args := m.Called(ctx, showtime)
	return args.Error(0)
}

func (m *MockShowtimeRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockShowtimeRepository) List(ctx context.Context, params pagination.Params, movieID, cinemaID uint, dateStr string) ([]models.Showtime, int64, error) {
	args := m.Called(ctx, params, movieID, cinemaID, dateStr)
	return args.Get(0).([]models.Showtime), args.Get(1).(int64), args.Error(2)
}

func (m *MockShowtimeRepository) DecrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error {
	args := m.Called(ctx, tx, showtimeID)
	return args.Error(0)
}

func (m *MockShowtimeRepository) IncrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error {
	args := m.Called(ctx, tx, showtimeID)
	return args.Error(0)
}

// MockTicketRepository
type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

func (m *MockTicketRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, ticket *models.Ticket) error {
	args := m.Called(ctx, tx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) GetByID(ctx context.Context, id uint) (*models.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetByShowtimeAndSeat(ctx context.Context, showtimeID uint, seatNumber string) (*models.Ticket, error) {
	args := m.Called(ctx, showtimeID, seatNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketRepository) ListByUser(ctx context.Context, userID uint, params pagination.Params) ([]models.Ticket, int64, error) {
	args := m.Called(ctx, userID, params)
	return args.Get(0).([]models.Ticket), args.Get(1).(int64), args.Error(2)
}

func (m *MockTicketRepository) ListAll(ctx context.Context, params pagination.Params, status string) ([]models.Ticket, int64, error) {
	args := m.Called(ctx, params, status)
	return args.Get(0).([]models.Ticket), args.Get(1).(int64), args.Error(2)
}

func (m *MockTicketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendTicketConfirmation(user *models.User, ticket *models.Ticket, showtime *models.Showtime) {
	m.Called(user, ticket, showtime)
}
