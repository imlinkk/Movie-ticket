package services

import (
	"context"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	"movie-ticket/internal/utils/pagination"
)

type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params) ([]models.User, int64, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, params pagination.Params) ([]models.User, int64, error) {
	return s.userRepo.List(ctx, params)
}
