package models

type SnippetView struct {
	SnippetID string  `gorm:"type:varchar(255);uniqueIndex:uq_snippet_views;not null" json:"-"`
	Hash      string  `gorm:"type:varchar(64);uniqueIndex:uq_snippet_views;not null" json:"-"`
	Snippet   Snippet `gorm:"foreignKey:SnippetID;references:SnippetID;constraint:OnDelete:CASCADE" json:"-"`
}

type GlobalViewsData struct {
	TotalViews int64 `json:"total_views"`
}

type GlobalViewsResponse struct {
	Success bool            `json:"success"`
	Data    GlobalViewsData `json:"data"`
}
