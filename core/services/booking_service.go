package services

import (
	"context"
	"fmt"
	"time"

	"movie-booking/api/v1/types"
	"movie-booking/config"
	"movie-booking/constants"
	"movie-booking/core/model"
	coretypes "movie-booking/core/types"
)

type bookingService struct {
	store model.DataStore
}

// NewBookingService creates a new booking service
func NewBookingService(clients *coretypes.Clients, store model.DataStore) BookingServiceInterface {
	return &bookingService{store: store}
}

// CreateBooking converts a locked seat into a confirmed booking
func (s *bookingService) CreateBooking(ctx context.Context, input *types.CreateBookingInput) (*types.BookingResponse, error) {
	// Check idempotency if key provided
	if input.IdempotencyKey != "" {
		existing, err := s.store.GetBookingByUserAndIdempotencyKey(ctx, input.UserID, input.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check idempotency: %w", err)
		}
		if existing != nil {
			// Return existing booking
			return &types.BookingResponse{
				BookingID: existing.ID,
				Status:    "CONFIRMED",
				Message:   "Booking already exists",
			}, nil
		}
	}

	// Begin transaction
	tx, err := s.store.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on panic or error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
	}()

	// Step 1: Lock the seat row using FOR UPDATE
	seat, err := tx.GetSeatByIDForUpdate(ctx, input.SeatID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	// Step 2: Validate
	lockDuration := config.GetSeatLockDuration()
	now := time.Now()

	// Must be LOCKED
	if seat.Status != string(constants.SeatStatusLocked) {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat is not locked")
	}

	// Lock must not be expired
	if seat.LockedAt == nil || now.Sub(*seat.LockedAt) > lockDuration {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat lock has expired")
	}

	// User must match (only locker can buy)
	if seat.UserID == nil || *seat.UserID != input.UserID {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat is locked by another user")
	}

	// Verify show matches
	if seat.ShowID != input.ShowID {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat does not belong to this show")
	}

	// Step 3: Update seat to SOLD
	updates := map[string]interface{}{
		"status":   string(constants.SeatStatusSold),
		"locked_at": nil,
		"user_id":   nil,
	}

	if err := tx.UpdateSeat(ctx, input.SeatID, updates); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to update seat: %w", err)
	}

	// Step 4: Create booking
	booking := &model.Booking{
		UserID:         input.UserID,
		ShowID:         input.ShowID,
		SeatID:         input.SeatID,
		IdempotencyKey: input.IdempotencyKey,
	}

	booking, err = tx.CreateBooking(ctx, booking)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Step 5: Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &types.BookingResponse{
		BookingID: booking.ID,
		Status:    "CONFIRMED",
		Message:   "Ticket sent to your email.",
	}, nil
}
