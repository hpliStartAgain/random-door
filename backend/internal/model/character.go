package model

import "time"

type Character struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CityID        int64     `gorm:"column:city_id;not null;index:idx_char_city;uniqueIndex:uk_char_city_name,priority:1" json:"city_id"`
	Name          string    `gorm:"column:name;type:varchar(128);not null;uniqueIndex:uk_char_city_name,priority:2" json:"name"`
	CharacterType string    `gorm:"column:character_type;type:varchar(32);not null" json:"character_type"`
	AvatarURL     *string   `gorm:"column:avatar_url;type:varchar(512)" json:"avatar_url,omitempty"`
	Persona       string    `gorm:"column:persona;type:text;not null" json:"-"`
	DialectStyle  *string   `gorm:"column:dialect_style;type:text" json:"dialect_style,omitempty"`
	Prompt        string    `gorm:"column:prompt;type:text;not null" json:"-"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Character) TableName() string { return "characters" }
