package services

import (
	"context"

	"movie-booking/api/v1/types"
	"movie-booking/core/model"
)

// AuthServiceInterface defines authentication operations
type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (*types.LoginResponse, error)
}

// MovieServiceInterface defines movie operations
type MovieServiceInterface interface {
	GetAllMovies(ctx context.Context) ([]model.Movie, error)
	GetMovieByID(ctx context.Context, id uint) (*model.Movie, error)
}

// ShowServiceInterface defines show operations
type ShowServiceInterface interface {
	GetShowsByMovieID(ctx context.Context, movieID uint) ([]model.Show, error)
	GetShowByID(ctx context.Context, id uint) (*model.Show, error)
}

// SeatServiceInterface defines seat operations
type SeatServiceInterface interface {
	GetSeatsByShowID(ctx context.Context, showID uint) ([]model.ShowSeat, error)
	LockSeat(ctx context.Context, seatID uint, userID uint) (*types.LockSeatResponse, error)
}

// BookingServiceInterface defines booking operations
type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, input *types.CreateBookingInput) (*types.BookingResponse, error)
}
