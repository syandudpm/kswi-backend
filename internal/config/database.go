package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

var db *gorm.DB

// InitDatabase initializes the database connection with GORM
func InitDatabase() error {
	log.Println("🔄 Initializing database connection...")

	// Get database configuration
	dbConfig := cfg.Database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		dbConfig.Host,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Port,
		dbConfig.SSLMode,
	)

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if IsDebug() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("✅ Database connected successfully to %s:%d", dbConfig.Host, dbConfig.Port)
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDatabase() first")
	}
	return db
}

// GetDatabaseDSN returns the formatted database DSN for GORM
func GetDatabaseDSN() string {
	dbConfig := cfg.Database
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		dbConfig.Host,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Port,
		dbConfig.SSLMode,
	)
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}

		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}

		log.Println("✅ Database connection closed")
	}
	return nil
}
