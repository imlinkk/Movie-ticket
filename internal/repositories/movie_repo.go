package repositories

import (
	"context"
	"errors"

	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"

	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *models.Movie) error
	GetByID(ctx context.Context, id uint) (*models.Movie, error)
	Update(ctx context.Context, movie *models.Movie) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, genre string) ([]models.Movie, int64, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Create(movie).Error
}

func (r *movieRepository) GetByID(ctx context.Context, id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.WithContext(ctx).First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("Movie not found", err)
		}
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Save(movie).Error
}

func (r *movieRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Movie{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return appErrors.NotFound("Movie not found to delete", nil)
	}
	return nil
}

func (r *movieRepository) List(ctx context.Context, params pagination.Params, genre string) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Movie{})
	if genre != "" {
		query = query.Where("genre ILIKE ?", "%"+genre+"%")
	}
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset()).
		Find(&movies).Error

	return movies, total, err
}
