package user

import (
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	db := config.GetDB() // Assumes this returns *gorm.DB
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", handler.CreateUser)
		// Later: userRoutes.GET("", handler.GetUsers)
		//        userRoutes.GET("/:id", handler.GetUserByID)
		//        etc.
	}
}
