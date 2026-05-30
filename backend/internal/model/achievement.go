package model

import "time"

type Achievement struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"column:code;type:varchar(64);not null;uniqueIndex:uk_ach_code" json:"code"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description *string   `gorm:"column:description;type:text" json:"description,omitempty"`
	RuleType    string    `gorm:"column:rule_type;type:varchar(64);not null" json:"rule_type"`
	RuleValue   string    `gorm:"column:rule_value;type:varchar(255);not null" json:"rule_value"`
	BadgeURL    *string   `gorm:"column:badge_url;type:varchar(512)" json:"badge_url,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Achievement) TableName() string { return "achievements" }
