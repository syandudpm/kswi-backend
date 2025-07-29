package router

import (
	"kswi-backend/internal/config"
	"kswi-backend/internal/modules/menu"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(gin.LoggerWithWriter(config.LogWriter()))
	r.Use(gin.Recovery())

	// CORS configuration
	corsConfig(r)

	// Common routes
	commonRoutes(r)

	// API routes
	api := r.Group("/api")
	{
		menu.RegisterRoutes(api)
	}

	return r
}

func corsConfig(r *gin.Engine) {
	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
}

func commonRoutes(r *gin.Engine) {
	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to KSWI Backend API!",
			"app":         config.Get().App.Name,
			"version":     config.Get().App.Version,
			"environment": config.Get().App.Environment,
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Check database health
		if err := config.HealthCheck(ctx); err != nil {
			config.GetSugaredLogger().Errorf("Health check failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"database":  "disconnected",
				"error":     err.Error(),
				"timestamp": time.Now().UTC(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"database":  "connected",
			"app":       config.GetAppInfo(),
			"timestamp": time.Now().UTC(),
		})
	})

	// API info endpoint
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "KSWI Backend API",
			"version": "v1",
			"endpoints": gin.H{
				"health": "/health",
				"menu":   "/api/menu",
			},
		})
	})
}
