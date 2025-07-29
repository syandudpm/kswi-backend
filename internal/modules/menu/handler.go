package menu

import (
	"net/http"
	"strconv"

	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// MenuTree returns the complete menu tree
func (h *Handler) MenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	logger := config.GetSugaredLogger()

	menuTree, err := h.svc.MenuTree(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch menu tree: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch menu tree",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    menuTree,
		"message": "Menu tree retrieved successfully",
	})
}

// GetActiveMenus returns all active menus (flat list)
func (h *Handler) GetActiveMenus(c *gin.Context) {
	ctx := c.Request.Context()
	logger := config.GetSugaredLogger()

	menus, err := h.svc.GetActiveMenus(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch active menus: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch active menus",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    menus,
		"message": "Active menus retrieved successfully",
	})
}

// GetMenuByID returns a specific menu by ID
func (h *Handler) GetMenuByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := config.GetSugaredLogger()

	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid menu ID",
			"message": "Menu ID must be a valid number",
		})
		return
	}

	menu, err := h.svc.GetMenuByID(ctx, uint(id))
	if err != nil {
		logger.Errorf("Failed to fetch menu by ID %d: %v", id, err)

		// Check if it's a "not found" error
		if err.Error() == "menu not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Menu not found",
				"message": "The requested menu does not exist",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch menu",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    menu,
		"message": "Menu retrieved successfully",
	})
}
