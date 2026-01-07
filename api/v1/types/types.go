package types

import (
	"net/http"
	"time"
)

// HandlerFunc is the signature for handler functions
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (*GenericAPIResponse, error)

// GenericAPIResponse is the standard API response format
type GenericAPIResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message,omitempty"`
	Values     interface{} `json:"values,omitempty"`
	Error      []FieldError `json:"error,omitempty"`
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LockSeatResponse represents the response for locking a seat
type LockSeatResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateBookingInput represents the input for creating a booking
type CreateBookingInput struct {
	ShowID         uint   `json:"show_id"`
	SeatID         uint   `json:"seat_id"`
	UserID         uint   // Set from JWT
	IdempotencyKey string // From header
}

// BookingResponse represents the response for a booking
type BookingResponse struct {
	BookingID uint   `json:"booking_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
