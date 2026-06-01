package models

import "time"

type RateLimit struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IP        string    `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	BlockedAt time.Time `gorm:"not null" json:"blocked_at"`
	UnblockAt time.Time `gorm:"not null;index" json:"unblock_at"`
	Reason    string    `gorm:"type:varchar(255);default:rate_limit" json:"reason"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
