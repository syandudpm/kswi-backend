package router

import (
	"kswi-backend/internal/config"
	"kswi-backend/internal/modules/menu"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	db := config.GetDB()

	r := gin.Default()

	corsConfig(r)
	commonRoutes(r, db)

	api := r.Group("/api")
	menu.RegisterRoutes(api)

	return r
}

func corsConfig(r *gin.Engine) {

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))
}

func commonRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		// if err := db.HealthCheck(c.Request.Context()); err != nil {
		// 	c.JSON(http.StatusServiceUnavailable, gin.H{
		// 		"status": "unhealthy",
		// 		"error":  err.Error(),
		// 	})
		// 	return
		// }
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
}
