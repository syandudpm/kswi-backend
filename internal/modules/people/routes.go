package people

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	repo := NewRepository(config.GetDB())
	svc := NewService(repo)
	handler := NewHandler(svc)

	peopleRoutes := r.Group("/people")
	{
		peopleRoutes.POST("/", handler.CreatePerson)
	}
}
