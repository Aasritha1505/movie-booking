package services

import (
	"context"
	"fmt"

	"movie-booking/core/model"
	coretypes "movie-booking/core/types"
)

type movieService struct {
	store model.DataStore
}

// NewMovieService creates a new movie service
func NewMovieService(clients *coretypes.Clients, store model.DataStore) MovieServiceInterface {
	return &movieService{store: store}
}

func (s *movieService) GetAllMovies(ctx context.Context) ([]model.Movie, error) {
	movies, err := s.store.GetAllMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	return movies, nil
}

func (s *movieService) GetMovieByID(ctx context.Context, id uint) (*model.Movie, error) {
	movie, err := s.store.GetMovieByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return movie, nil
}
