package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"movie-booking/api/v1/types"
	"github.com/gorilla/mux"
)

// ValidateAndParseLoginRequest parses and validates login request
func ValidateAndParseLoginRequest(r *http.Request) (*types.LoginRequest, error) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &req, nil
}

// ValidateAndParseBookingRequest parses and validates booking request
func ValidateAndParseBookingRequest(r *http.Request, userID uint) (*types.CreateBookingInput, error) {
	var req struct {
		ShowID uint `json:"show_id"`
		SeatID uint `json:"seat_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if req.ShowID == 0 {
		return nil, fmt.Errorf("show_id is required")
	}
	if req.SeatID == 0 {
		return nil, fmt.Errorf("seat_id is required")
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	return &types.CreateBookingInput{
		ShowID:         req.ShowID,
		SeatID:         req.SeatID,
		UserID:         userID,
		IdempotencyKey: idempotencyKey,
	}, nil
}

// ParseUintFromPath extracts a uint from URL path variable
func ParseUintFromPath(r *http.Request, key string) (uint, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[key]
	if !ok {
		return 0, fmt.Errorf("%s not found in path", key)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return uint(id), nil
}

// ParseMovieIDFromPath extracts movie ID from path
func ParseMovieIDFromPath(r *http.Request) (uint, error) {
	return ParseUintFromPath(r, "id")
}

// ParseShowIDFromPath extracts show ID from path
func ParseShowIDFromPath(r *http.Request) (uint, error) {
	return ParseUintFromPath(r, "id")
}

// ParseSeatIDFromPath extracts seat ID from path
func ParseSeatIDFromPath(r *http.Request) (uint, error) {
	return ParseUintFromPath(r, "id")
}
