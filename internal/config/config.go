package config

import "fmt"

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
}

var cfg *Config

// Get returns the global config instance
func Get() *Config {
	return cfg
}

// IsProduction returns true if the environment is production
func IsProduction() bool {
	return cfg.App.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func IsDevelopment() bool {
	return cfg.App.Environment == "development"
}

// IsDebug returns true if debug mode is enabled
func IsDebug() bool {
	return cfg.App.Debug
}

// GetServerAddress returns the full server address
func GetServerAddress() string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}

// GetAppInfo returns formatted application information
func GetAppInfo() string {
	return fmt.Sprintf("%s v%s (%s)", cfg.App.Name, cfg.App.Version, cfg.App.Environment)
}
