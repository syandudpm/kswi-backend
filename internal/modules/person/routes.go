package person

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	repo := NewRepository(config.GetDB())
	svc := NewService(repo)
	handler := NewHandler(svc)

	personRoutes := r.Group("/persons")
	{
		personRoutes.POST("/", handler.CreatePerson)
	}
}
