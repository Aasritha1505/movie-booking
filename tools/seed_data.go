package main

import (
	"fmt"
	"log"
	"time"

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

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database. Seeding data...")

	// Create movies
	movies := []struct {
		Title       string
		Description string
		Duration    int
		Rating      string
	}{
		{
			Title:       "The Matrix",
			Description: "A computer hacker learns about the true nature of reality and his role in the war against its controllers.",
			Duration:    136,
			Rating:      "R",
		},
		{
			Title:       "Inception",
			Description: "A skilled thief is given a chance at redemption if he can accomplish the impossible task of inception.",
			Duration:    148,
			Rating:      "PG-13",
		},
		{
			Title:       "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			Duration:    169,
			Rating:      "PG-13",
		},
	}

	movieIDs := []uint{}
	for _, m := range movies {
		var count int64
		db.Model(&struct {
			ID    uint   `gorm:"primaryKey"`
			Title string `gorm:"column:title"`
		}{}).Table("movies").Where("title = ?", m.Title).Count(&count)

		if count == 0 {
			result := db.Exec(`
				INSERT INTO movies (title, description, duration_mins, content_rating, created_at, updated_at)
				VALUES (?, ?, ?, ?, NOW(), NOW())
			`, m.Title, m.Description, m.Duration, m.Rating)
			if result.Error != nil {
				log.Printf("Error inserting movie %s: %v", m.Title, result.Error)
				continue
			}
			fmt.Printf("✓ Created movie: %s\n", m.Title)
		} else {
			fmt.Printf("- Movie already exists: %s\n", m.Title)
		}

		// Get the movie ID
		var movieID uint
		db.Raw("SELECT id FROM movies WHERE title = ?", m.Title).Scan(&movieID)
		movieIDs = append(movieIDs, movieID)
	}

	// Create theatres
	theatres := []struct {
		Name     string
		Location string
	}{
		{
			Name:     "PVR Cinemas",
			Location: "Downtown Mall",
		},
		{
			Name:     "IMAX Theatre",
			Location: "City Center",
		},
		{
			Name:     "Cineplex",
			Location: "Shopping Plaza",
		},
	}

	theatreIDs := []uint{}
	for _, t := range theatres {
		var count int64
		db.Model(&struct {
			ID   uint   `gorm:"primaryKey"`
			Name string `gorm:"column:name"`
		}{}).Table("theatres").Where("name = ? AND location = ?", t.Name, t.Location).Count(&count)

		if count == 0 {
			result := db.Exec(`
				INSERT INTO theatres (name, location, created_at, updated_at)
				VALUES (?, ?, NOW(), NOW())
			`, t.Name, t.Location)
			if result.Error != nil {
				log.Printf("Error inserting theatre %s: %v", t.Name, result.Error)
				continue
			}
			fmt.Printf("✓ Created theatre: %s - %s\n", t.Name, t.Location)
		} else {
			fmt.Printf("- Theatre already exists: %s - %s\n", t.Name, t.Location)
		}

		// Get the theatre ID
		var theatreID uint
		db.Raw("SELECT id FROM theatres WHERE name = ? AND location = ?", t.Name, t.Location).Scan(&theatreID)
		theatreIDs = append(theatreIDs, theatreID)
	}

	// Create shows (one show per movie per theatre, at different times)
	now := time.Now()
	showTimes := []time.Time{
		now.Add(2 * time.Hour),                    // 2 hours from now
		now.Add(5 * time.Hour),                    // 5 hours from now
		now.AddDate(0, 0, 1).Add(2 * time.Hour),  // Tomorrow at 2 PM
		now.AddDate(0, 0, 1).Add(6 * time.Hour),  // Tomorrow at 6 PM
	}

	showIDs := []uint{}
	showIndex := 0
	for _, movieID := range movieIDs {
		for _, theatreID := range theatreIDs {
			if showIndex < len(showTimes) {
				showTime := showTimes[showIndex%len(showTimes)]
				showIndex++

				// Check if show exists and get its ID
				var existingShowID uint
				var count int64
				db.Raw("SELECT id FROM shows WHERE movie_id = ? AND theatre_id = ? AND ABS(TIMESTAMPDIFF(SECOND, start_time, ?)) < 60",
					movieID, theatreID, showTime).Scan(&existingShowID)
				if existingShowID > 0 {
					count = 1
				}

				var showID uint
				if count == 0 {
					// Create new show
					result := db.Exec(`
						INSERT INTO shows (movie_id, theatre_id, start_time, created_at, updated_at)
						VALUES (?, ?, ?, NOW(), NOW())
					`, movieID, theatreID, showTime)
					if result.Error != nil {
						log.Printf("Error inserting show: %v", result.Error)
						continue
					}
					fmt.Printf("✓ Created show: Movie ID %d at Theatre ID %d at %s\n",
						movieID, theatreID, showTime.Format("2006-01-02 15:04:05"))
					
					// Get the newly created show ID
					db.Raw("SELECT id FROM shows WHERE movie_id = ? AND theatre_id = ? AND ABS(TIMESTAMPDIFF(SECOND, start_time, ?)) < 60",
						movieID, theatreID, showTime).Scan(&showID)
				} else {
					// Use existing show ID
					showID = existingShowID
					fmt.Printf("- Show already exists: Movie ID %d at Theatre ID %d (ID: %d)\n", movieID, theatreID, showID)
				}

				if showID > 0 {
					showIDs = append(showIDs, showID)
				}
			}
		}
	}

	// Create seats for each show (50 seats: A1-A10, B1-B10, C1-C10, D1-D10, E1-E10)
	rows := []string{"A", "B", "C", "D", "E"}
	seatsPerRow := 10

	seatsCreated := 0
	for _, showID := range showIDs {
		// Check if seats already exist for this show
		var seatCount int64
		db.Model(&struct {
			ID uint `gorm:"primaryKey"`
		}{}).Table("show_seats").Where("show_id = ?", showID).Count(&seatCount)

		if seatCount == 0 {
			// Create seats for this show
			for _, row := range rows {
				for i := 1; i <= seatsPerRow; i++ {
					seatName := fmt.Sprintf("%s%d", row, i)
					result := db.Exec(`
						INSERT INTO show_seats (show_id, seat_name, status, created_at, updated_at)
						VALUES (?, ?, 'AVAILABLE', NOW(), NOW())
					`, showID, seatName)
					if result.Error != nil {
						log.Printf("Error inserting seat %s for show %d: %v", seatName, showID, result.Error)
						continue
					}
					seatsCreated++
				}
			}
			fmt.Printf("✓ Created 50 seats for show ID %d\n", showID)
		} else {
			fmt.Printf("- Seats already exist for show ID %d (%d seats)\n", showID, seatCount)
		}
	}

	fmt.Printf("\n✅ Seeding complete!\n")
	fmt.Printf("   - Movies: %d\n", len(movieIDs))
	fmt.Printf("   - Theatres: %d\n", len(theatreIDs))
	fmt.Printf("   - Shows: %d\n", len(showIDs))
	fmt.Printf("   - New seats created: %d\n", seatsCreated)
	fmt.Printf("   - Total seats: %d\n", len(showIDs)*50)
}
