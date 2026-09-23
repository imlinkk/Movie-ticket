package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	"movie-ticket/internal/utils/cache"
	"movie-ticket/internal/utils/pagination"
)

type MovieService interface {
	Create(ctx context.Context, req *models.CreateMovieRequest) (*models.Movie, error)
	GetByID(ctx context.Context, id uint) (*models.Movie, error)
	Update(ctx context.Context, id uint, req *models.UpdateMovieRequest) (*models.Movie, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params, genre string) ([]models.Movie, int64, error)
}

type movieService struct {
	movieRepo repositories.MovieRepository
	cache     *cache.RedisClient
}

func NewMovieService(movieRepo repositories.MovieRepository, cache *cache.RedisClient) MovieService {
	return &movieService{
		movieRepo: movieRepo,
		cache:     cache,
	}
}

func (s *movieService) Create(ctx context.Context, req *models.CreateMovieRequest) (*models.Movie, error) {
	movie := &models.Movie{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Genre:           req.Genre,
		ReleaseDate:     req.ReleaseDate,
		PosterURL:       req.PosterURL,
		TrailerURL:      req.TrailerURL,
		Rating:          req.Rating,
	}

	if err := s.movieRepo.Create(ctx, movie); err != nil {
		return nil, err
	}

	// Invalidate movie list caches
	if s.cache != nil {
		_ = s.cache.DeleteByPrefix(ctx, "movies:list:")
	}

	return movie, nil
}

func (s *movieService) GetByID(ctx context.Context, id uint) (*models.Movie, error) {
	cacheKey := fmt.Sprintf("movies:item:%d", id)

	// Try reading from cache
	if s.cache != nil {
		var cachedMovie models.Movie
		if err := s.cache.Get(ctx, cacheKey, &cachedMovie); err == nil {
			log.Printf("[Redis Cache HIT] Movie ID: %d", id)
			return &cachedMovie, nil
		}
	}

	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache for 10 minutes
	if s.cache != nil {
		_ = s.cache.Set(ctx, cacheKey, movie, 10*time.Minute)
	}

	return movie, nil
}

func (s *movieService) Update(ctx context.Context, id uint, req *models.UpdateMovieRequest) (*models.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.DurationMinutes > 0 {
		movie.DurationMinutes = req.DurationMinutes
	}
	if req.Genre != "" {
		movie.Genre = req.Genre
	}
	if req.ReleaseDate != nil {
		movie.ReleaseDate = *req.ReleaseDate
	}
	if req.PosterURL != "" {
		movie.PosterURL = req.PosterURL
	}
	if req.TrailerURL != "" {
		movie.TrailerURL = req.TrailerURL
	}
	if req.Rating != nil {
		movie.Rating = *req.Rating
	}

	if err := s.movieRepo.Update(ctx, movie); err != nil {
		return nil, err
	}

	// Invalidate caches
	if s.cache != nil {
		_ = s.cache.Delete(ctx, fmt.Sprintf("movies:item:%d", id))
		_ = s.cache.DeleteByPrefix(ctx, "movies:list:")
	}

	return movie, nil
}

func (s *movieService) Delete(ctx context.Context, id uint) error {
	if err := s.movieRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate caches
	if s.cache != nil {
		_ = s.cache.Delete(ctx, fmt.Sprintf("movies:item:%d", id))
		_ = s.cache.DeleteByPrefix(ctx, "movies:list:")
	}

	return nil
}

func (s *movieService) List(ctx context.Context, params pagination.Params, genre string) ([]models.Movie, int64, error) {
	cacheKey := fmt.Sprintf("movies:list:page:%d:limit:%d:sort:%s:order:%s:genre:%s:search:%s",
		params.Page, params.Limit, params.SortBy, params.Order, genre, params.Search)

	type listCacheData struct {
		Movies []models.Movie
		Total  int64
	}

	if s.cache != nil {
		var cached listCacheData
		if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
			log.Printf("[Redis Cache HIT] Movies list: %s", cacheKey)
			return cached.Movies, cached.Total, nil
		}
	}

	movies, total, err := s.movieRepo.List(ctx, params, genre)
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, cacheKey, listCacheData{Movies: movies, Total: total}, 5*time.Minute)
	}

	return movies, total, nil
}
