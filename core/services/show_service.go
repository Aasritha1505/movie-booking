package services

import (
	"context"
	"fmt"

	"movie-booking/core/model"
	coretypes "movie-booking/core/types"
)

type showService struct {
	store model.DataStore
}

// NewShowService creates a new show service
func NewShowService(clients *coretypes.Clients, store model.DataStore) ShowServiceInterface {
	return &showService{store: store}
}

func (s *showService) GetShowsByMovieID(ctx context.Context, movieID uint) ([]model.Show, error) {
	shows, err := s.store.GetShowsByMovieID(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shows: %w", err)
	}
	return shows, nil
}

func (s *showService) GetShowByID(ctx context.Context, id uint) (*model.Show, error) {
	show, err := s.store.GetShowByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get show: %w", err)
	}
	return show, nil
}
