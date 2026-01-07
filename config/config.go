package config

import (
	"time"

	"github.com/spf13/viper"
)

var settings *viper.Viper

// Init initializes the configuration
func Init() error {
	settings = viper.New()
	settings.SetConfigType("env")
	settings.AutomaticEnv()

	// Set defaults
	settings.SetDefault("DATABASE_HOST", "localhost")
	settings.SetDefault("DATABASE_PORT", "3306")
	settings.SetDefault("DATABASE_USER", "root")
	settings.SetDefault("DATABASE_PASSWORD", "password")
	settings.SetDefault("DATABASE_NAME", "movie_booking")
	settings.SetDefault("SERVER_PORT", "8080")
	settings.SetDefault("HANDLER_TIMEOUT", "30s")
	settings.SetDefault("JWT_SECRET", "change-me-in-production")
	settings.SetDefault("JWT_EXPIRY", "15m")
	settings.SetDefault("SEAT_LOCK_DURATION", "10m")

	return nil
}

// GetSettings returns the viper instance
func GetSettings() *viper.Viper {
	return settings
}

// Database configuration
func GetDatabaseHost() string {
	return settings.GetString("DATABASE_HOST")
}

func GetDatabasePort() string {
	return settings.GetString("DATABASE_PORT")
}

func GetDatabaseUser() string {
	return settings.GetString("DATABASE_USER")
}

func GetDatabasePassword() string {
	return settings.GetString("DATABASE_PASSWORD")
}

func GetDatabaseName() string {
	return settings.GetString("DATABASE_NAME")
}

// Server configuration
func GetServerPort() string {
	return settings.GetString("SERVER_PORT")
}

func GetHandlerTimeout() time.Duration {
	return settings.GetDuration("HANDLER_TIMEOUT")
}

// JWT configuration
func GetJWTSecret() string {
	return settings.GetString("JWT_SECRET")
}

func GetJWTExpiry() time.Duration {
	return settings.GetDuration("JWT_EXPIRY")
}

// Seat lock configuration
func GetSeatLockDuration() time.Duration {
	return settings.GetDuration("SEAT_LOCK_DURATION")
}
