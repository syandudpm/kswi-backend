package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateUser(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		case ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}
