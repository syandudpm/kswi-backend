package config

import (
	"fmt"
	"log"
)

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// InitApp initializes the entire application
func InitApp() error {
	log.Println("🚀 Starting application initialization...")

	// 1. Initialize Viper configuration
	if err := InitViper(); err != nil {
		return fmt.Errorf("failed to initialize viper: %w", err)
	}

	// 2. Initialize Logger
	if err := InitLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 3. Initialize Database
	if err := InitDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 4. Initialize Redis
	if err := InitRedis(); err != nil {
		return fmt.Errorf("failed to initialize redis: %w", err)
	}

	// 5. Initialize JWT
	if err := InitJWT(); err != nil {
		return fmt.Errorf("failed to initialize JWT: %w", err)
	}

	GetLogger().Infof("✅ Application initialized successfully: %s", GetAppInfo())
	return nil
}

// ShutdownApp gracefully shuts down the application
func ShutdownApp() error {
	log.Println("🛑 Shutting down application...")

	var errors []error

	// Close database connection
	if err := CloseDatabase(); err != nil {
		errors = append(errors, err)
		GetLogger().Errorf("Failed to close database: %v", err)
	}

	// Close Redis connection
	if err := CloseRedis(); err != nil {
		errors = append(errors, err)
		GetLogger().Errorf("Failed to close Redis: %v", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	log.Println("✅ Application shutdown completed")
	return nil
}
