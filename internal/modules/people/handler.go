package people

import (
	"kswi-backend/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreatePerson(c *gin.Context) {
	var req CreatePersonRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.HandleValidationError(err))
		return
	}

	// Call service to create people
	people, err := h.service.CreatePerson(&req)
	if err != nil {
		c.Error(err)
		return
	}

	// Return success response
	c.JSON(201, gin.H{
		"success": true,
		"message": "Person created successfully",
		"data":    people,
	})
}
