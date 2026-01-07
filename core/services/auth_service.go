package services

import (
	"context"
	"fmt"
	"time"

	"movie-booking/api/v1/types"
	"movie-booking/config"
	"movie-booking/core/model"
	coretypes "movie-booking/core/types"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	store model.DataStore
}

// NewAuthService creates a new auth service
func NewAuthService(clients *coretypes.Clients, store model.DataStore) AuthServiceInterface {
	return &authService{store: store}
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, email, password string) (*types.LoginResponse, error) {
	// Get user by email
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &types.LoginResponse{
		Token: token,
		User: types.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// generateJWT creates a JWT token with RS256 (simplified to HS256 for MVP, but structure supports RS256)
func (s *authService) generateJWT(userID uint, email string) (string, error) {
	expiry := config.GetJWTExpiry()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.GetJWTSecret()
	return token.SignedString([]byte(secret))
}
