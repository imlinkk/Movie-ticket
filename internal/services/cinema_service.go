package services

import (
	"context"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	"movie-ticket/internal/utils/pagination"
)

type CinemaService interface {
	Create(ctx context.Context, req *models.CreateCinemaRequest) (*models.Cinema, error)
	GetByID(ctx context.Context, id uint) (*models.Cinema, error)
	Update(ctx context.Context, id uint, req *models.UpdateCinemaRequest) (*models.Cinema, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, city string) ([]models.Cinema, int64, error)
}

type cinemaService struct {
	cinemaRepo repositories.CinemaRepository
}

func NewCinemaService(cinemaRepo repositories.CinemaRepository) CinemaService {
	return &cinemaService{cinemaRepo: cinemaRepo}
}

func (s *cinemaService) Create(ctx context.Context, req *models.CreateCinemaRequest) (*models.Cinema, error) {
	cinema := &models.Cinema{
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		TotalRooms: req.TotalRooms,
	}

	if err := s.cinemaRepo.Create(ctx, cinema); err != nil {
		return nil, err
	}
	return cinema, nil
}

func (s *cinemaService) GetByID(ctx context.Context, id uint) (*models.Cinema, error) {
	return s.cinemaRepo.GetByID(ctx, id)
}

func (s *cinemaService) Update(ctx context.Context, id uint, req *models.UpdateCinemaRequest) (*models.Cinema, error) {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		cinema.Name = req.Name
	}
	if req.Address != "" {
		cinema.Address = req.Address
	}
	if req.City != "" {
		cinema.City = req.City
	}
	if req.TotalRooms > 0 {
		cinema.TotalRooms = req.TotalRooms
	}

	if err := s.cinemaRepo.Update(ctx, cinema); err != nil {
		return nil, err
	}

	return cinema, nil
}

func (s *cinemaService) Delete(ctx context.Context, id uint) error {
	return s.cinemaRepo.Delete(ctx, id)
}

func (s *cinemaService) List(ctx context.Context, params pagination.Params, city string) ([]models.Cinema, int64, error) {
	return s.cinemaRepo.List(ctx, params, city)
}
