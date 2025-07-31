package pagination

import (
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	CurrentPage   int `json:"current_page"`
	PerPage       int `json:"per_page"`
	TotalPages    int `json:"total_pages"`
	Total         int `json:"total"`
	TotalFiltered int `json:"total_filtered"`
}

type PaginationLinks struct {
	First string  `json:"first"`
	Last  string  `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

type ResponseParam struct {
	Ctx      *gin.Context
	Req      PaginationRequest
	Data     any
	Total    int
	Filtered int
}

type PaginationResponse struct {
	Success bool                   `json:"success"`
	Data    PaginationResponseData `json:"data"`
}

type PaginationResponseData struct {
	Data  any             `json:"data"`
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
}

func BuildMeta(currentPage, perPage, total, totalFiltered int) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalFiltered) / float64(perPage)))

	return PaginationMeta{
		CurrentPage:   currentPage,
		PerPage:       perPage,
		Total:         total,
		TotalFiltered: totalFiltered,
		TotalPages:    totalPages,
	}
}

func BuildLinks(c *gin.Context, req PaginationRequest, totalPages int) PaginationLinks {
	var prevPage, nextPage *string

	basePath := c.Request.URL.Path

	queryParams := c.Request.URL.Query()
	queryParams.Del("page")
	extraQueryParams := queryParams.Encode()

	buildURL := func(page int) string {
		url := fmt.Sprintf("%s?page=%d", basePath, page)
		if extraQueryParams != "" {
			url += "&" + extraQueryParams
		}
		return url
	}

	if req.Page > 1 {
		prev := buildURL(req.Page - 1)
		prevPage = &prev
	}

	if req.Page < totalPages {
		next := buildURL(req.Page + 1)
		nextPage = &next
	}

	first := buildURL(1)
	last := buildURL(totalPages)

	return PaginationLinks{
		First: first,
		Last:  last,
		Prev:  prevPage,
		Next:  nextPage,
	}
}

func BuildResponse(d ResponseParam) *PaginationResponse {

	meta := BuildMeta(d.Req.Page, d.Req.PerPage, d.Total, d.Filtered)
	links := BuildLinks(d.Ctx, d.Req, meta.TotalPages)

	return &PaginationResponse{
		Success: true,
		Data: PaginationResponseData{
			Data:  d.Data,
			Meta:  meta,
			Links: links,
		},
	}
}
