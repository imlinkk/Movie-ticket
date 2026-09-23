package services

import (
	"context"
	"time"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/pagination"
)

type ShowtimeService interface {
	Create(ctx context.Context, req *models.CreateShowtimeRequest) (*models.Showtime, error)
	GetByID(ctx context.Context, id uint) (*models.Showtime, error)
	Update(ctx context.Context, id uint, req *models.UpdateShowtimeRequest) (*models.Showtime, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, movieID, cinemaID uint, dateStr string) ([]models.Showtime, int64, error)
}

type showtimeService struct {
	showtimeRepo repositories.ShowtimeRepository
	movieRepo    repositories.MovieRepository
	cinemaRepo   repositories.CinemaRepository
}

func NewShowtimeService(
	showtimeRepo repositories.ShowtimeRepository,
	movieRepo repositories.MovieRepository,
	cinemaRepo repositories.CinemaRepository,
) ShowtimeService {
	return &showtimeService{
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		cinemaRepo:   cinemaRepo,
	}
}

func (s *showtimeService) Create(ctx context.Context, req *models.CreateShowtimeRequest) (*models.Showtime, error) {
	movie, err := s.movieRepo.GetByID(ctx, req.MovieID)
	if err != nil {
		return nil, appErrors.NotFound("Movie not found for showtime", err)
	}

	_, err = s.cinemaRepo.GetByID(ctx, req.CinemaID)
	if err != nil {
		return nil, appErrors.NotFound("Cinema not found for showtime", err)
	}

	endTime := req.StartTime.Add(time.Duration(movie.DurationMinutes) * time.Minute)

	showtime := &models.Showtime{
		MovieID:        req.MovieID,
		CinemaID:       req.CinemaID,
		RoomName:       req.RoomName,
		StartTime:      req.StartTime,
		EndTime:        endTime,
		Price:          req.Price,
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.TotalSeats,
	}

	if err := s.showtimeRepo.Create(ctx, showtime); err != nil {
		return nil, err
	}

	return s.showtimeRepo.GetByID(ctx, showtime.ID)
}

func (s *showtimeService) GetByID(ctx context.Context, id uint) (*models.Showtime, error) {
	return s.showtimeRepo.GetByID(ctx, id)
}

func (s *showtimeService) Update(ctx context.Context, id uint, req *models.UpdateShowtimeRequest) (*models.Showtime, error) {
	showtime, err := s.showtimeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.RoomName != "" {
		showtime.RoomName = req.RoomName
	}
	if req.StartTime != nil {
		showtime.StartTime = *req.StartTime
		if showtime.Movie.DurationMinutes > 0 {
			showtime.EndTime = req.StartTime.Add(time.Duration(showtime.Movie.DurationMinutes) * time.Minute)
		}
	}
	if req.Price != nil {
		showtime.Price = *req.Price
	}

	if err := s.showtimeRepo.Update(ctx, showtime); err != nil {
		return nil, err
	}

	return showtime, nil
}

func (s *showtimeService) Delete(ctx context.Context, id uint) error {
	return s.showtimeRepo.Delete(ctx, id)
}

func (s *showtimeService) List(ctx context.Context, params pagination.Params, movieID, cinemaID uint, dateStr string) ([]models.Showtime, int64, error) {
	return s.showtimeRepo.List(ctx, params, movieID, cinemaID, dateStr)
}
