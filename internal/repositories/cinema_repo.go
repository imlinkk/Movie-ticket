package repositories

import (
	"context"
	"errors"

	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"

	"gorm.io/gorm"
)

type CinemaRepository interface {
	Create(ctx context.Context, cinema *models.Cinema) error
	GetByID(ctx context.Context, id uint) (*models.Cinema, error)
	Update(ctx context.Context, cinema *models.Cinema) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, city string) ([]models.Cinema, int64, error)
}

type cinemaRepository struct {
	db *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) CinemaRepository {
	return &cinemaRepository{db: db}
}

func (r *cinemaRepository) Create(ctx context.Context, cinema *models.Cinema) error {
	return r.db.WithContext(ctx).Create(cinema).Error
}

func (r *cinemaRepository) GetByID(ctx context.Context, id uint) (*models.Cinema, error) {
	var cinema models.Cinema
	if err := r.db.WithContext(ctx).First(&cinema, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("Cinema not found", err)
		}
		return nil, err
	}
	return &cinema, nil
}

func (r *cinemaRepository) Update(ctx context.Context, cinema *models.Cinema) error {
	return r.db.WithContext(ctx).Save(cinema).Error
}

func (r *cinemaRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Cinema{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return appErrors.NotFound("Cinema not found to delete", nil)
	}
	return nil
}

func (r *cinemaRepository) List(ctx context.Context, params pagination.Params, city string) ([]models.Cinema, int64, error) {
	var cinemas []models.Cinema
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Cinema{})
	if city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR address ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset()).
		Find(&cinemas).Error

	return cinemas, total, err
}
