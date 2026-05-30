package model

import "time"

type ChatMessage struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"column:user_id;not null;index:idx_cm_user_char_time" json:"user_id"`
	CityID      int64     `gorm:"column:city_id;not null;index:idx_cm_city" json:"city_id"`
	CharacterID int64     `gorm:"column:character_id;not null;index:idx_cm_user_char_time;index:idx_cm_character" json:"character_id"`
	Role        string    `gorm:"column:role;type:varchar(16);not null" json:"role"`
	Content     string    `gorm:"column:content;type:text;not null" json:"content"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index:idx_cm_user_char_time" json:"created_at"`
}

func (ChatMessage) TableName() string { return "chat_messages" }
