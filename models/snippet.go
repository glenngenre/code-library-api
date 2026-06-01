package models

import "time"

type Snippet struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SnippetID   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"snippet_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Code        string    `gorm:"type:text;not null" json:"code"`
	Language    string    `gorm:"type:varchar(50);not null;index" json:"language"`
	Filename    string    `gorm:"type:varchar(255);not null" json:"filename"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(100);index" json:"category"`
	IsPublic    bool      `gorm:"default:true;not null" json:"is_public"`
	IsFeatured  bool      `gorm:"default:false;not null" json:"is_featured"`
	Views       int       `gorm:"default:0;not null" json:"views"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_created_at,sort:desc" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateSnippetRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Code        string `json:"code" binding:"required,min=5"`
	Language    string `json:"language" binding:"required,min=1,max=50"`
	Filename    string `json:"filename" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsPublic    bool   `json:"is_public"`
}

type SnippetResponse struct {
	Success bool    `json:"success"`
	Data    Snippet `json:"data"`
}

type CursorPayload struct {
	CreatedAt time.Time `json:"ca"`
	SnippetID string    `json:"sid"`
}

type SnippetCursorListResponse struct {
	Success bool       `json:"success"`
	Data    []Snippet  `json:"data"`
	Meta    CursorMeta `json:"meta"`
}
