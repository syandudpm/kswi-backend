package menu

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	// Initialize dependencies with proper dependency injection
	repo := NewRepository(config.GetDB()) // Pass database connection
	svc := NewService(repo)               // Service takes Repository interface
	h := NewHandler(svc)                  // Handler takes Service interface

	// Create menu route group
	menuRoutes := r.Group("/menu")
	{
		// GET Routes - Read Operations

		// Get menu tree (hierarchical structure)
		menuRoutes.GET("", h.MenuTree)

	}
}
