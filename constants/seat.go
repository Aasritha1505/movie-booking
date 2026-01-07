package constants

// SeatStatus represents the status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusLocked    SeatStatus = "LOCKED"
	SeatStatusSold      SeatStatus = "SOLD"
)

// ValidSeatStatuses returns all valid seat statuses
var ValidSeatStatuses = []SeatStatus{
	SeatStatusAvailable,
	SeatStatusLocked,
	SeatStatusSold,
}
