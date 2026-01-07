package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run create_user.go <email> <password> <name>")
		fmt.Println("Example: go run create_user.go test@example.com password123 'Test User'")
		os.Exit(1)
	}

	email := os.Args[1]
	password := os.Args[2]
	name := os.Args[3]

	// Generate password hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	dsn := "root:password@tcp(localhost:3306)/movie_booking"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Insert user
	query := "INSERT INTO users (email, password_hash, name, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	result, err := db.Exec(query, email, string(hash), name)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	id, _ := result.LastInsertId()
	fmt.Printf("✓ User created successfully!\n")
	fmt.Printf("  ID: %d\n", id)
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Password: %s\n", password)
}
