package oss

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	repo := NewRepository(config.GetDB())
	svc := NewService(repo)
	h := NewHandler(svc)

	routes := r.Group("/oss")
	{
		routes.GET("/tree", h.Test)
		routes.POST("/dt", h.DtDatabase)
	}
}
