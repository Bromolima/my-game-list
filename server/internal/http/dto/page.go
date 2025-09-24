package dto

type PageResponse[T any] struct {
	Data       []T   `query:"data"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
}
