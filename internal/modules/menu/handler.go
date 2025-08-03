package menu

import (
	"kswi-backend/internal/shared/api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetMenuTree godoc
// @Summary Get active menu tree
// @Description Retrieves hierarchical menu structure with active menus only
// @Tags menu
// @Produce json
// @Success 200 {object} api.APIResponse{data=[]MenuResponse}
// @Failure 500 {object} api.APIResponse
// @Router /api/menu/tree [get]
func (h *Handler) GetMenuTree(c *gin.Context) {
	ctx := c.Request.Context()

	menus, err := h.service.GetMenuTree(ctx)
	if err != nil {
		// Error will be handled by error middleware
		_ = c.Error(err)
		return
	}

	response := api.APIResponse{
		Success: true,
		Message: "Menu tree retrieved successfully",
		Data:    menus,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetMenuByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	ctx := c.Request.Context()
	menu, err := h.service.GetMenuByID(ctx, uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if menu == nil {
		c.JSON(404, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(200, menu)
}
