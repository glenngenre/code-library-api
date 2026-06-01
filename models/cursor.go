package models

type CursorMeta struct {
	PageSize   int     `json:"page_size"`
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
	HasMore    bool    `json:"has_more"`
}
