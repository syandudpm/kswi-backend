package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// InitViper initializes Viper configuration
func InitViper() error {
	v := viper.New()

	// Set config file details
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add multiple config paths
	v.AddConfigPath("./config")          // ./config/config.yaml
	v.AddConfigPath("./internal/config") // ./internal/config/config.yaml
	v.AddConfigPath(".")                 // ./config.yaml
	v.AddConfigPath("/etc/kswi")         // /etc/kswi/config.yaml (for production)

	// Enable environment variable support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("KSWI")

	// Set default values
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Unmarshal config
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "KSWI Backend")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.idle_timeout", 60)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "kswi_db")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 300)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conn", 5)

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.access_token_ttl", 3600)    // 1 hour
	v.SetDefault("jwt.refresh_token_ttl", 604800) // 7 days
	v.SetDefault("jwt.issuer", "kswi-backend")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
}
