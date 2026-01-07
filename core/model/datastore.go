package model

import (
	"context"
)

//go:generate mockgen -destination=../../datastore/fake/fake.go -package=fake movie-booking/core/model DataStore

// DataStore is the main interface for data access
type DataStore interface {
	// Composed interfaces
	UserStore
	MovieStore
	ShowStore
	ShowSeatStore
	BookingStore

	// Transaction support
	Begin(ctx context.Context) (DataStore, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	InTransactionMode(ctx context.Context) bool
	WithContext(ctx context.Context) DataStore
}

// UserStore handles user operations
type UserStore interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// MovieStore handles movie operations
type MovieStore interface {
	GetAllMovies(ctx context.Context) ([]Movie, error)
	GetMovieByID(ctx context.Context, id uint) (*Movie, error)
}

// ShowStore handles show operations
type ShowStore interface {
	GetShowsByMovieID(ctx context.Context, movieID uint) ([]Show, error)
	GetShowByID(ctx context.Context, id uint) (*Show, error)
}

// ShowSeatStore handles seat operations
type ShowSeatStore interface {
	GetSeatsByShowID(ctx context.Context, showID uint) ([]ShowSeat, error)
	GetSeatByIDForUpdate(ctx context.Context, id uint) (*ShowSeat, error) // FOR UPDATE lock
	UpdateSeat(ctx context.Context, id uint, updates map[string]interface{}) error
	CreateSeat(ctx context.Context, seat *ShowSeat) (*ShowSeat, error)
}

// BookingStore handles booking operations
type BookingStore interface {
	CreateBooking(ctx context.Context, booking *Booking) (*Booking, error)
	GetBookingByUserAndIdempotencyKey(ctx context.Context, userID uint, idempotencyKey string) (*Booking, error)
	GetBookingByID(ctx context.Context, id uint) (*Booking, error)
}
