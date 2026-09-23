package repositories

import (
	"context"
	"errors"
	"time"

	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"

	"gorm.io/gorm"
)

type ShowtimeRepository interface {
	Create(ctx context.Context, showtime *models.Showtime) error
	GetByID(ctx context.Context, id uint) (*models.Showtime, error)
	Update(ctx context.Context, showtime *models.Showtime) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, movieID, cinemaID uint, dateStr string) ([]models.Showtime, int64, error)
	DecrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error
	IncrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error
}

type showtimeRepository struct {
	db *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) ShowtimeRepository {
	return &showtimeRepository{db: db}
}

func (r *showtimeRepository) Create(ctx context.Context, showtime *models.Showtime) error {
	return r.db.WithContext(ctx).Create(showtime).Error
}

func (r *showtimeRepository) GetByID(ctx context.Context, id uint) (*models.Showtime, error) {
	var showtime models.Showtime
	if err := r.db.WithContext(ctx).Preload("Movie").Preload("Cinema").First(&showtime, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("Showtime not found", err)
		}
		return nil, err
	}
	return &showtime, nil
}

func (r *showtimeRepository) Update(ctx context.Context, showtime *models.Showtime) error {
	return r.db.WithContext(ctx).Save(showtime).Error
}

func (r *showtimeRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Showtime{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return appErrors.NotFound("Showtime not found to delete", nil)
	}
	return nil
}

func (r *showtimeRepository) List(ctx context.Context, params pagination.Params, movieID, cinemaID uint, dateStr string) ([]models.Showtime, int64, error) {
	var showtimes []models.Showtime
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Showtime{}).Preload("Movie").Preload("Cinema")

	if movieID > 0 {
		query = query.Where("movie_id = ?", movieID)
	}
	if cinemaID > 0 {
		query = query.Where("cinema_id = ?", cinemaID)
	}
	if dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset()).
		Find(&showtimes).Error

	return showtimes, total, err
}

func (r *showtimeRepository) DecrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	res := db.WithContext(ctx).Model(&models.Showtime{}).
		Where("id = ? AND available_seats > 0", showtimeID).
		Update("available_seats", gorm.Expr("available_seats - 1"))

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return appErrors.Conflict("No available seats left for this showtime", nil)
	}
	return nil
}

func (r *showtimeRepository) IncrementAvailableSeats(ctx context.Context, tx *gorm.DB, showtimeID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Model(&models.Showtime{}).
		Where("id = ?", showtimeID).
		Update("available_seats", gorm.Expr("available_seats + 1")).Error
}
