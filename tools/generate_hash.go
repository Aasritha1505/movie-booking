package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", string(hash))
	
	// Verify it works
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
	} else {
		fmt.Println("✓ Hash verification successful!")
	}
}
