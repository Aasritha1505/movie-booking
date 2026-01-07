package main

import (
	"fmt"
	"log"

	"movie-booking/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Initialize config
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetDatabaseUser(),
		config.GetDatabasePassword(),
		config.GetDatabaseHost(),
		config.GetDatabasePort(),
		config.GetDatabaseName(),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("=== MOVIE BOOKING DATABASE DATA ===\n")

	// Count and show movies
	var movieCount int64
	db.Table("movies").Count(&movieCount)
	fmt.Printf("📽️  MOVIES (%d total)\n", movieCount)
	fmt.Println("-----------------------------------")
	
	type Movie struct {
		ID          uint   `gorm:"column:id"`
		Title       string `gorm:"column:title"`
		Description string `gorm:"column:description"`
		Duration    int    `gorm:"column:duration_mins"`
		Rating      string `gorm:"column:content_rating"`
	}
	var movies []Movie
	db.Table("movies").Find(&movies)
	for _, m := range movies {
		fmt.Printf("  ID: %d | %s (%d mins, %s)\n", m.ID, m.Title, m.Duration, m.Rating)
	}
	fmt.Println()

	// Count and show theatres
	var theatreCount int64
	db.Table("theatres").Count(&theatreCount)
	fmt.Printf("🎭 THEATRES (%d total)\n", theatreCount)
	fmt.Println("-----------------------------------")
	
	type Theatre struct {
		ID       uint   `gorm:"column:id"`
		Name     string `gorm:"column:name"`
		Location string `gorm:"column:location"`
	}
	var theatres []Theatre
	db.Table("theatres").Find(&theatres)
	for _, t := range theatres {
		fmt.Printf("  ID: %d | %s - %s\n", t.ID, t.Name, t.Location)
	}
	fmt.Println()

	// Count and show shows
	var showCount int64
	db.Table("shows").Count(&showCount)
	fmt.Printf("🎬 SHOWS (%d total)\n", showCount)
	fmt.Println("-----------------------------------")
	
	type Show struct {
		ID        uint   `gorm:"column:id"`
		MovieID   uint   `gorm:"column:movie_id"`
		TheatreID uint   `gorm:"column:theatre_id"`
		StartTime string `gorm:"column:start_time"`
	}
	var shows []Show
	db.Table("shows").Order("start_time").Find(&shows)
	for _, s := range shows {
		var movieTitle, theatreName string
		db.Table("movies").Where("id = ?", s.MovieID).Select("title").Scan(&movieTitle)
		db.Table("theatres").Where("id = ?", s.TheatreID).Select("name").Scan(&theatreName)
		fmt.Printf("  ID: %d | %s @ %s | %s\n", s.ID, movieTitle, theatreName, s.StartTime)
	}
	fmt.Println()

	// Count and show seats
	var seatCount int64
	db.Table("show_seats").Count(&seatCount)
	fmt.Printf("🪑 SEATS (%d total)\n", seatCount)
	fmt.Println("-----------------------------------")
	
	type Seat struct {
		ID       uint   `gorm:"column:id"`
		ShowID   uint   `gorm:"column:show_id"`
		SeatName string `gorm:"column:seat_name"`
		Status   string `gorm:"column:status"`
	}
	
	// Group by show
	for _, s := range shows {
		var seats []Seat
		db.Table("show_seats").Where("show_id = ?", s.ID).Find(&seats)
		available := 0
		locked := 0
		sold := 0
		for _, seat := range seats {
			switch seat.Status {
			case "AVAILABLE":
				available++
			case "LOCKED":
				locked++
			case "SOLD":
				sold++
			}
		}
		var movieTitle string
		db.Table("movies").Where("id = ?", s.MovieID).Select("title").Scan(&movieTitle)
		fmt.Printf("  Show ID %d (%s): %d seats - %d available, %d locked, %d sold\n",
			s.ID, movieTitle, len(seats), available, locked, sold)
	}
	fmt.Println()

	// Count and show users
	var userCount int64
	db.Table("users").Count(&userCount)
	fmt.Printf("👤 USERS (%d total)\n", userCount)
	fmt.Println("-----------------------------------")
	
	type User struct {
		ID    uint   `gorm:"column:id"`
		Email string `gorm:"column:email"`
		Name  string `gorm:"column:name"`
	}
	var users []User
	db.Table("users").Find(&users)
	for _, u := range users {
		fmt.Printf("  ID: %d | %s (%s)\n", u.ID, u.Email, u.Name)
	}
	fmt.Println()

	// Count and show bookings
	var bookingCount int64
	db.Table("bookings").Count(&bookingCount)
	fmt.Printf("🎫 BOOKINGS (%d total)\n", bookingCount)
	fmt.Println("-----------------------------------")
	
	type Booking struct {
		ID     uint   `gorm:"column:id"`
		UserID uint   `gorm:"column:user_id"`
		ShowID uint   `gorm:"column:show_id"`
		SeatID uint   `gorm:"column:seat_id"`
	}
	var bookings []Booking
	db.Table("bookings").Find(&bookings)
	for _, b := range bookings {
		var userEmail, seatName string
		db.Table("users").Where("id = ?", b.UserID).Select("email").Scan(&userEmail)
		db.Table("show_seats").Where("id = ?", b.SeatID).Select("seat_name").Scan(&seatName)
		fmt.Printf("  ID: %d | User: %s | Show: %d | Seat: %s\n", b.ID, userEmail, b.ShowID, seatName)
	}
	fmt.Println()
}
