package menu

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	repo := NewRepository()
	svc := NewService(repo)
	h := NewHandler(svc)

	routes := r.Group("/menu")
	{
		// Get menu tree (hierarchical structure)
		routes.GET("", h.MenuTree)
		routes.GET("/tree", h.MenuTree) // Alternative endpoint for clarity

		// Get all active menus (flat list)
		routes.GET("/active", h.GetActiveMenus)

		// Get specific menu by ID
		routes.GET("/:id", h.GetMenuByID)
	}
}
