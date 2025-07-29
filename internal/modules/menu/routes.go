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
		routes.GET("", h.MenuTree)
	}
}
