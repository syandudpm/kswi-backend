package oss

import (
	"kswi-backend/internal/shared/errors"
	"kswi-backend/internal/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Test(c *gin.Context) {}

func (h *Handler) DtDatabase(c *gin.Context) {
	var req DtDatabaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.HandleValidationError(err))
		return
	}

	data, total, filtered, err := h.svc.DtDatabase(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pagination.BuildResponse(
		pagination.ResponseParam{
			Ctx:      c,
			Req:      req.PaginationRequest,
			Data:     data,
			Total:    total,
			Filtered: filtered,
		},
	))

}
