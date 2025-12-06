package utils

import (
	"math"

	"gorm.io/gorm"
)

type PaginationResult struct {
	Results      interface{} `json:"results"`
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	TotalPages   int         `json:"totalPages"`
	TotalResults int64       `json:"totalResults"`
}

type PaginationScope struct {
	Page  int
	Limit int
	Sort  string
}

// Paginate returns a GORM scope for pagination
func (p *PaginationScope) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	}
}

// GetPaginationResult calculates metadata (totalPages, etc.)
func GetPaginationResult(totalRows int64, page, limit int, data interface{}) PaginationResult {
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	return PaginationResult{
		Results:      data,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalResults: totalRows,
	}
}