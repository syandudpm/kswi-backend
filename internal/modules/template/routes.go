package template

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	repo := NewRepository(config.GetDB())
	svc := NewService(repo)
	handler := NewHandler(svc)

	menuRoutes := r.Group("/template")
	{
		menuRoutes.GET("/tree", handler.Test)
	}
}
