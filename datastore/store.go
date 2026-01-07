package datastore

import (
	"context"
	"errors"
	"fmt"

	"movie-booking/core/model"
	"gorm.io/gorm"
)

// DBStore implements the DataStore interface using GORM
type DBStore struct {
	db                *gorm.DB
	inTransactionMode bool
}

// NewDataStore creates a new DataStore instance
func NewDataStore(db *gorm.DB) model.DataStore {
	return &DBStore{db: db, inTransactionMode: false}
}

// Begin starts a transaction
func (ds *DBStore) Begin(ctx context.Context) (model.DataStore, error) {
	tx := ds.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return &DBStore{db: tx, inTransactionMode: true}, nil
}

// Commit commits the transaction
func (ds *DBStore) Commit(ctx context.Context) error {
	if !ds.inTransactionMode {
		return errors.New("not in transaction mode")
	}
	if err := ds.db.WithContext(ctx).Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the transaction
func (ds *DBStore) Rollback(ctx context.Context) error {
	if !ds.inTransactionMode {
		return errors.New("not in transaction mode")
	}
	if err := ds.db.WithContext(ctx).Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

// InTransactionMode returns whether we're in a transaction
func (ds *DBStore) InTransactionMode(ctx context.Context) bool {
	return ds.inTransactionMode
}

// WithContext returns a new store with the given context
func (ds *DBStore) WithContext(ctx context.Context) model.DataStore {
	return &DBStore{db: ds.db.WithContext(ctx), inTransactionMode: ds.inTransactionMode}
}

// UserStore implementation

func (ds *DBStore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := ds.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (ds *DBStore) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := ds.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (ds *DBStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := ds.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// MovieStore implementation

func (ds *DBStore) GetAllMovies(ctx context.Context) ([]model.Movie, error) {
	var movies []model.Movie
	if err := ds.db.WithContext(ctx).Find(&movies).Error; err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	return movies, nil
}

func (ds *DBStore) GetMovieByID(ctx context.Context, id uint) (*model.Movie, error) {
	var movie model.Movie
	if err := ds.db.WithContext(ctx).Where("id = ?", id).First(&movie).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return &movie, nil
}

// ShowStore implementation

func (ds *DBStore) GetShowsByMovieID(ctx context.Context, movieID uint) ([]model.Show, error) {
	var shows []model.Show
	if err := ds.db.WithContext(ctx).
		Preload("Movie").
		Preload("Theatre").
		Where("movie_id = ?", movieID).
		Find(&shows).Error; err != nil {
		return nil, fmt.Errorf("failed to get shows: %w", err)
	}
	return shows, nil
}

func (ds *DBStore) GetShowByID(ctx context.Context, id uint) (*model.Show, error) {
	var show model.Show
	if err := ds.db.WithContext(ctx).
		Preload("Movie").
		Preload("Theatre").
		Where("id = ?", id).
		First(&show).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("show not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get show: %w", err)
	}
	return &show, nil
}

// ShowSeatStore implementation

func (ds *DBStore) GetSeatsByShowID(ctx context.Context, showID uint) ([]model.ShowSeat, error) {
	var seats []model.ShowSeat
	if err := ds.db.WithContext(ctx).
		Where("show_id = ?", showID).
		Order("seat_name").
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}
	return seats, nil
}

// GetSeatByIDForUpdate locks the seat row using FOR UPDATE
func (ds *DBStore) GetSeatByIDForUpdate(ctx context.Context, id uint) (*model.ShowSeat, error) {
	var seat model.ShowSeat
	// Use raw SQL with FOR UPDATE for row-level locking
	// GORM doesn't have a direct Locking clause in v1.26, so we use Set("gorm:query_option", "FOR UPDATE")
	if err := ds.db.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", id).
		First(&seat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("seat not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get seat for update: %w", err)
	}
	return &seat, nil
}

func (ds *DBStore) UpdateSeat(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := ds.db.WithContext(ctx).
		Model(&model.ShowSeat{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update seat: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("seat not found or no changes made")
	}
	return nil
}

func (ds *DBStore) CreateSeat(ctx context.Context, seat *model.ShowSeat) (*model.ShowSeat, error) {
	if err := ds.db.WithContext(ctx).Create(seat).Error; err != nil {
		return nil, fmt.Errorf("failed to create seat: %w", err)
	}
	return seat, nil
}

// BookingStore implementation

func (ds *DBStore) CreateBooking(ctx context.Context, booking *model.Booking) (*model.Booking, error) {
	if err := ds.db.WithContext(ctx).Create(booking).Error; err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}
	return booking, nil
}

func (ds *DBStore) GetBookingByUserAndIdempotencyKey(ctx context.Context, userID uint, idempotencyKey string) (*model.Booking, error) {
	var booking model.Booking
	if err := ds.db.WithContext(ctx).
		Where("user_id = ? AND idempotency_key = ?", userID, idempotencyKey).
		First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is OK for idempotency check
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return &booking, nil
}

func (ds *DBStore) GetBookingByID(ctx context.Context, id uint) (*model.Booking, error) {
	var booking model.Booking
	if err := ds.db.WithContext(ctx).
		Preload("User").
		Preload("Show").
		Preload("Seat").
		Where("id = ?", id).
		First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("booking not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return &booking, nil
}
