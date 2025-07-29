package menu

import "github.com/gin-gonic/gin"

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) MenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menuTree, err := h.svc.MenuTree(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch menu tree"})
		return
	}
	c.JSON(200, menuTree)
}
