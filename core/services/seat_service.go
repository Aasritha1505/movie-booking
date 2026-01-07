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

type seatService struct {
	store model.DataStore
}

// NewSeatService creates a new seat service
func NewSeatService(clients *coretypes.Clients, store model.DataStore) SeatServiceInterface {
	return &seatService{store: store}
}

func (s *seatService) GetSeatsByShowID(ctx context.Context, showID uint) ([]model.ShowSeat, error) {
	seats, err := s.store.GetSeatsByShowID(ctx, showID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	// Lazy lock expiration: treat expired locks as AVAILABLE
	lockDuration := config.GetSeatLockDuration()
	now := time.Now()
	for i := range seats {
		if seats[i].Status == string(constants.SeatStatusLocked) && seats[i].LockedAt != nil {
			if now.Sub(*seats[i].LockedAt) > lockDuration {
				// Lock expired, treat as available (optionally update in DB)
				seats[i].Status = string(constants.SeatStatusAvailable)
				seats[i].LockedAt = nil
				seats[i].UserID = nil
			}
		}
	}

	return seats, nil
}

// LockSeat implements the core concurrency strategy with row-level locking
func (s *seatService) LockSeat(ctx context.Context, seatID uint, userID uint) (*types.LockSeatResponse, error) {
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
	seat, err := tx.GetSeatByIDForUpdate(ctx, seatID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	// Step 2: Validate - allow if AVAILABLE or if LOCKED but expired
	lockDuration := config.GetSeatLockDuration()
	now := time.Now()
	canLock := false

	if seat.Status == string(constants.SeatStatusAvailable) {
		canLock = true
	} else if seat.Status == string(constants.SeatStatusLocked) && seat.LockedAt != nil {
		// Check if lock is expired
		if now.Sub(*seat.LockedAt) > lockDuration {
			canLock = true
		} else {
			// Lock is still valid
			tx.Rollback(ctx)
			return nil, fmt.Errorf("seat is already locked")
		}
	} else if seat.Status == string(constants.SeatStatusSold) {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat is already sold")
	}

	if !canLock {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("seat cannot be locked")
	}

	// Step 3: Update seat to LOCKED
	lockedAt := now
	updates := map[string]interface{}{
		"status":   string(constants.SeatStatusLocked),
		"locked_at": lockedAt,
		"user_id":   userID,
	}

	if err := tx.UpdateSeat(ctx, seatID, updates); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to lock seat: %w", err)
	}

	// Step 4: Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	expiresAt := lockedAt.Add(lockDuration)
	return &types.LockSeatResponse{
		Message:   "Locked",
		ExpiresAt: expiresAt,
	}, nil
}
