package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kswi-backend/internal/config"
	"kswi-backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application
	if err := config.InitApp(); err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Get logger
	logger := config.GetSugaredLogger()

	// Set Gin mode based on environment
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router with all routes
	ginRouter := router.SetupRouter()

	// Create HTTP server
	server := &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      ginRouter,
		ReadTimeout:  time.Duration(config.Get().Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Get().Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Get().Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("🚀 Server starting on %s", config.GetServerAddress())
		logger.Infof("📊 Application: %s", config.GetAppInfo())
		logger.Infof("🌍 Environment: %s", config.Get().App.Environment)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	// Shutdown application components
	if err := config.ShutdownApp(); err != nil {
		logger.Errorf("Error during application shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("✅ Server exited successfully")
}
