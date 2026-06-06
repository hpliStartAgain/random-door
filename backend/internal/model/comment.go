package model

import "time"

type Comment struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TargetType string    `gorm:"column:target_type;type:varchar(32);not null;index:idx_comments_target_time,priority:1" json:"target_type"`
	TargetID   int64     `gorm:"column:target_id;not null;index:idx_comments_target_time,priority:2" json:"target_id"`
	UserID     *int64    `gorm:"column:user_id;index:idx_comments_user_time,priority:1" json:"user_id,omitempty"`
	Nickname   string    `gorm:"column:nickname;type:varchar(64);not null" json:"nickname"`
	Content    string    `gorm:"column:content;type:varchar(500);not null" json:"content"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_comments_target_time,priority:3;index:idx_comments_user_time,priority:2" json:"created_at"`
}

func (Comment) TableName() string { return "comments" }
