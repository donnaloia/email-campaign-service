package models

// PaginationParams for incoming requests
type PaginationParams struct {
	Page     int
	PageSize int
}

// PaginatedResponse is a generic type for paginated results
type PaginatedResponse[T any] struct {
	Results     []T `json:"results"`
	Total       int `json:"total"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](results []T, total, currentPage, pageSize int) *PaginatedResponse[T] {
	totalPages := (total + pageSize - 1) / pageSize
	return &PaginatedResponse[T]{
		Results:     results,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}
}
