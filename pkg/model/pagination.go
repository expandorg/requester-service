package model

// Pagination model
type Pagination struct {
	Current uint64 `json:"current"`
	Total   uint64 `json:"total"`
}

// NewPagination model
func NewPagination(current uint64, total uint64) *Pagination {
	return &Pagination{
		Current: current,
		Total:   total,
	}
}
