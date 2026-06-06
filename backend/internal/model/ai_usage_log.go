package model

import "time"

type AIUsageLog struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;uniqueIndex:uk_ai_usage_user_type_date,priority:1" json:"user_id"`
	UsageType string    `gorm:"column:usage_type;type:varchar(32);not null;uniqueIndex:uk_ai_usage_user_type_date,priority:2" json:"usage_type"`
	UsageDate time.Time `gorm:"column:usage_date;type:date;not null;uniqueIndex:uk_ai_usage_user_type_date,priority:3;index:idx_ai_usage_date" json:"usage_date"`
	Count     int       `gorm:"column:count;not null;default:0" json:"count"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (AIUsageLog) TableName() string { return "ai_usage_logs" }
