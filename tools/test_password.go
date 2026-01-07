package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get the stored hash from database
	storedHash := "$2a$10$9pw4SqNSMv4U0q892DJPr.gD4n90ZUCCNIsTnv20oeomq/9q3TWVK"
	password := "password123"
	
	fmt.Printf("Testing password: %s\n", password)
	fmt.Printf("Stored hash: %s\n", storedHash)
	
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Printf("✗ Password verification FAILED: %v\n", err)
	} else {
		fmt.Println("✓ Password verification SUCCESSFUL!")
	}
}
