package utils

import (
	"math"
	
	"github.com/galex-do/test-machine/internal/models"
)

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, total int) models.PaginationResponse {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25 // default page size
	}
	if pageSize > 100 {
		pageSize = 100 // max page size
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	return models.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

// GetOffsetAndLimit calculates SQL OFFSET and LIMIT values
func GetOffsetAndLimit(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset = (page - 1) * pageSize
	limit = pageSize

	return offset, limit
}

// DefaultPagination returns default pagination request
func DefaultPagination() models.PaginationRequest {
	return models.PaginationRequest{
		Page:     1,
		PageSize: 25,
	}
}