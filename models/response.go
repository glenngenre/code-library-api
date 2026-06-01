package models

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type EmptyResponse struct {
	Success bool `json:"success"`
}
