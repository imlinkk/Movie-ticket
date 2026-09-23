package tests

import (
	"context"
	"testing"
	"time"

	"movie-ticket/internal/models"
	"movie-ticket/internal/services"
	"movie-ticket/internal/utils/pagination"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMovieService_Create_Success(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	movieService := services.NewMovieService(mockRepo, nil)

	ctx := context.Background()
	req := &models.CreateMovieRequest{
		Title:           "Inception",
		Description:     "A thief who steals corporate secrets through the use of dream-sharing technology.",
		DurationMinutes: 148,
		Genre:           "Sci-Fi",
		ReleaseDate:     time.Now(),
		Rating:          8.8,
	}

	mockRepo.On("Create", ctx, mock.MatchedBy(func(m *models.Movie) bool {
		return m.Title == "Inception" && m.Genre == "Sci-Fi"
	})).Return(nil)

	movie, err := movieService.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, "Inception", movie.Title)
	assert.Equal(t, 148, movie.DurationMinutes)
	mockRepo.AssertExpectations(t)
}

func TestMovieService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	movieService := services.NewMovieService(mockRepo, nil)

	ctx := context.Background()
	expectedMovie := &models.Movie{
		ID:              1,
		Title:           "Interstellar",
		DurationMinutes: 169,
		Genre:           "Sci-Fi",
		Rating:          8.7,
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(expectedMovie, nil)

	movie, err := movieService.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, uint(1), movie.ID)
	assert.Equal(t, "Interstellar", movie.Title)
	mockRepo.AssertExpectations(t)
}

func TestMovieService_List_Success(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	movieService := services.NewMovieService(mockRepo, nil)

	ctx := context.Background()
	params := pagination.Params{
		Page:   1,
		Limit:  10,
		SortBy: "created_at",
		Order:  "desc",
	}

	expectedMovies := []models.Movie{
		{ID: 1, Title: "Movie 1", Genre: "Action"},
		{ID: 2, Title: "Movie 2", Genre: "Action"},
	}

	mockRepo.On("List", ctx, params, "Action").Return(expectedMovies, int64(2), nil)

	movies, total, err := movieService.List(ctx, params, "Action")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, movies, 2)
	assert.Equal(t, "Movie 1", movies[0].Title)
	mockRepo.AssertExpectations(t)
}
