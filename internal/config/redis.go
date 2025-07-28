package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Password    string `mapstructure:"password"`
	Database    int    `mapstructure:"database"`
	MaxRetries  int    `mapstructure:"max_retries"`
	PoolSize    int    `mapstructure:"pool_size"`
	MinIdleConn int    `mapstructure:"min_idle_conn"`
}

var rdb *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() error {
	log.Println("🔄 Initializing Redis connection...")

	// Get Redis configuration
	redisConfig := cfg.Redis

	// Create Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.Database,
		MaxRetries:   redisConfig.MaxRetries,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConn,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("✅ Redis connected successfully to %s:%d", redisConfig.Host, redisConfig.Port)
	return nil
}

// GetRedis returns the Redis client
func GetRedis() *redis.Client {
	if rdb == nil {
		log.Fatal("Redis not initialized. Call InitRedis() first")
	}
	return rdb
}

// GetRedisAddress returns the formatted Redis address
func GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			return fmt.Errorf("failed to close redis: %w", err)
		}
		log.Println("✅ Redis connection closed")
	}
	return nil
}
