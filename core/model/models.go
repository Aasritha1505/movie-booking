package model

import (
	"time"
)

// Movie represents a movie entity
type Movie struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Title        string `gorm:"type:varchar(500);not null" json:"title"`
	Description  string `gorm:"type:text" json:"description"`
	DurationMins int    `json:"duration_mins"`
	ContentRating string `gorm:"type:varchar(50)" json:"rating"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Movie) TableName() string {
	return "movies"
}

// Theatre represents a theatre entity
type Theatre struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Location string `gorm:"type:varchar(255)" json:"location"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Theatre) TableName() string {
	return "theatres"
}

// Show represents a movie show at a theatre
type Show struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MovieID   uint      `gorm:"not null;index" json:"movie_id"`
	TheatreID uint      `gorm:"not null;index" json:"theatre_id"`
	StartTime time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	
	// Relations
	Movie   Movie   `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Theatre Theatre `gorm:"foreignKey:TheatreID" json:"theatre,omitempty"`
}

func (Show) TableName() string {
	return "shows"
}

// User represents a user entity
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// ShowSeat represents a seat for a specific show
type ShowSeat struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	ShowID   uint       `gorm:"not null;index" json:"show_id"`
	SeatName string     `gorm:"type:varchar(10);not null" json:"seat_name"`
	Status   string     `gorm:"type:varchar(50);default:'AVAILABLE';index" json:"status"` // AVAILABLE, LOCKED, SOLD
	LockedAt *time.Time  `gorm:"type:timestamp NULL" json:"locked_at,omitempty"`
	UserID   *uint       `gorm:"index" json:"user_id,omitempty"` // WHO locked this seat
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	
	// Relations
	Show Show `gorm:"foreignKey:ShowID" json:"show,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ShowSeat) TableName() string {
	return "show_seats"
}

// Booking represents a confirmed booking
type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ShowID    uint      `gorm:"not null;index" json:"show_id"`
	SeatID    uint      `gorm:"not null;index" json:"seat_id"`
	IdempotencyKey string `gorm:"type:varchar(255);index" json:"-"` // For idempotency
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	
	// Relations
	User User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Show Show     `gorm:"foreignKey:ShowID" json:"show,omitempty"`
	Seat ShowSeat `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}
