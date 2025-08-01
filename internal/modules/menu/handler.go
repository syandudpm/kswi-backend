package menu

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetMenuTree godoc
// @Summary Get active menu tree
// @Tags menu
// @Produce json
// @Success 200 {array} MenuResponse
// @Router /menu/tree [get]
func (h *Handler) GetMenuTree(c *gin.Context) {
	tree, err := h.GetMenuTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load menu"})
		return
	}
	c.JSON(http.StatusOK, tree)
}

// CreateMenu godoc
// @Summary Create a new menu
// @Tags menu
// @Accept json
// @Produce json
// @Param input body CreateMenuInput true "Menu data"
// @Success 201
// @Router /menu [post]
func (h *Handler) CreateMenu(c *gin.Context) {
	var input CreateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.CreateMenu(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu"})
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateMenu godoc
// @Summary Update a menu
// @Tags menu
// @Accept json
// @Produce json
// @Param id path uint true "Menu ID"
// @Param input body UpdateMenuInput true "Update data"
// @Success 204
// @Router /menu/{id} [patch]
func (h *Handler) UpdateMenu(c *gin.Context) {
	var input UpdateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var menuID uint
	_, err := fmt.Sscanf(id, "%d", &menuID)
	if err != nil || menuID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	err = h.UpdateMenu(c.Request.Context(), menuID, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteMenu godoc
// @Summary Delete a menu (soft delete)
// @Tags menu
// @Success 204
// @Router /menu/{id} [delete]
func (h *Handler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	var menuID uint
	_, err := fmt.Sscanf(id, "%d", &menuID)
	if err != nil || menuID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	err = h.DeleteMenu(c.Request.Context(), menuID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
