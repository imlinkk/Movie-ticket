package pagination

import (
	"math"
	"strconv"

	"movie-ticket/internal/models"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page    int
	Limit   int
	SortBy  string
	Order   string
	Search  string
}

func GetPaginationParams(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	search := c.Query("search")

	return Params{
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

func BuildPagination(page, limit int, totalItems int64) models.Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return models.Pagination{
		CurrentPage: page,
		PageSize:    limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}
