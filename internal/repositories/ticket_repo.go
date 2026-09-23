package repositories

import (
	"context"
	"errors"

	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetDB() *gorm.DB
	CreateWithTx(ctx context.Context, tx *gorm.DB, ticket *models.Ticket) error
	GetByID(ctx context.Context, id uint) (*models.Ticket, error)
	GetByShowtimeAndSeat(ctx context.Context, showtimeID uint, seatNumber string) (*models.Ticket, error)
	ListByUser(ctx context.Context, userID uint, params pagination.Params) ([]models.Ticket, int64, error)
	ListAll(ctx context.Context, params pagination.Params, status string) ([]models.Ticket, int64, error)
	Update(ctx context.Context, ticket *models.Ticket) error
	Delete(ctx context.Context, id uint) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *ticketRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, ticket *models.Ticket) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return appErrors.Conflict("This seat has already been booked for this showtime", err)
		}
		return err
	}
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Cinema").
		First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("Ticket not found", err)
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByShowtimeAndSeat(ctx context.Context, showtimeID uint, seatNumber string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(ctx).
		Where("showtime_id = ? AND seat_number = ? AND status != ?", showtimeID, seatNumber, models.TicketStatusCancelled).
		First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) ListByUser(ctx context.Context, userID uint, params pagination.Params) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Ticket{}).
		Where("user_id = ?", userID).
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Cinema")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset()).
		Find(&tickets).Error

	return tickets, total, err
}

func (r *ticketRepository) ListAll(ctx context.Context, params pagination.Params, status string) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Ticket{}).
		Preload("User").
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Cinema")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset()).
		Find(&tickets).Error

	return tickets, total, err
}

func (r *ticketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Ticket{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return appErrors.NotFound("Ticket not found to delete", nil)
	}
	return nil
}
