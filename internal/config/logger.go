package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig holds logging-specific configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

var logger *zap.Logger
var sugar *zap.SugaredLogger

// InitLogger initializes the application logger
func InitLogger() error {
	log.Println("🔄 Initializing logger...")

	// Get log configuration
	logConfig := cfg.Log

	// Parse log level
	var level zapcore.Level
	switch strings.ToLower(logConfig.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		return fmt.Errorf("invalid log level: %s", logConfig.Level)
	}

	// Configure encoder
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	switch strings.ToLower(logConfig.Format) {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "text":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return fmt.Errorf("unsupported log format: %s", logConfig.Format)
	}

	// Configure output
	var writeSyncer zapcore.WriteSyncer
	switch strings.ToLower(logConfig.Output) {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "file":
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writeSyncer = zapcore.AddSync(file)
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core and logger
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// Add application context to all logs
	logger = logger.With(
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Environment),
	)

	// Create sugared logger for easier use
	sugar = logger.Sugar()

	log.Printf("✅ Logger initialized with level=%s format=%s", logConfig.Level, logConfig.Format)

	// Test log
	logger.Info("Logger initialized successfully")

	return nil
}

// GetLogger returns the structured zap logger
func GetLogger() *zap.Logger {
	if logger == nil {
		log.Fatal("Logger not initialized. Call InitLogger() first")
	}
	return logger
}

// GetSugaredLogger returns the sugared zap logger (easier API)
func GetSugaredLogger() *zap.SugaredLogger {
	if sugar == nil {
		log.Fatal("Logger not initialized. Call InitLogger() first")
	}
	return sugar
}

// zapWriter implements io.Writer for zap logger
type zapWriter struct {
	logger *zap.Logger
}

func (w *zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(strings.TrimSpace(string(p)))
	return len(p), nil
}

// LogWriter returns an io.Writer for the logger (useful for HTTP middleware)
func LogWriter() io.Writer {
	return &zapWriter{logger: GetLogger()}
}

// SetLogOutput sets the log output destination
func SetLogOutput(output io.Writer) {
	if logger != nil {
		// Create new core with the new output
		core := logger.Core()
		newCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(output),
			core,
		)
		logger = zap.New(newCore, zap.AddCaller(), zap.AddCallerSkip(1))
		sugar = logger.Sugar()
	}
}

// Sync flushes any buffered log entries (call before application exits)
func Sync() {
	if logger != nil {
		logger.Sync()
	}
	if sugar != nil {
		sugar.Sync()
	}
}
