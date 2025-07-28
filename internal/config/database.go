package config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

var db *gorm.DB

// InitDatabase initializes the database connection with GORM
func InitDatabase() error {
	GetSugaredLogger().Info("🔄 Initializing database connection...")

	// Get database configuration
	dbConfig := cfg.Database

	// Create MySQL DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.Charset,
		dbConfig.ParseTime,
		dbConfig.Loc,
	)

	// Configure GORM logger based on environment
	var gormLogger gormlogger.Interface
	if IsDebug() {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Info)
	} else {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Silent)
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
		GetSugaredLogger().Errorf("Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		GetSugaredLogger().Errorf("Failed to get underlying sql.DB: %v", err)
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
		GetSugaredLogger().Errorf("Failed to ping database: %v", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	GetSugaredLogger().Infof("✅ Database connected successfully to %s:%d", dbConfig.Host, dbConfig.Port)
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if db == nil {
		GetSugaredLogger().Fatal("Database not initialized. Call InitDatabase() first")
	}
	return db
}

// GetDatabaseDSN returns the formatted database DSN for GORM
func GetDatabaseDSN() string {
	dbConfig := cfg.Database
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.Charset,
		dbConfig.ParseTime,
		dbConfig.Loc,
	)
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			GetSugaredLogger().Errorf("Failed to get underlying sql.DB: %v", err)
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}

		if err := sqlDB.Close(); err != nil {
			GetSugaredLogger().Errorf("Failed to close database: %v", err)
			return fmt.Errorf("failed to close database: %w", err)
		}

		GetSugaredLogger().Info("✅ Database connection closed")
	}
	return nil
}
