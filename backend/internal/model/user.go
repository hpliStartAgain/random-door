package model

import "time"

type User struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"user_id"`
	AnonymousID   string    `gorm:"column:anonymous_id;type:varchar(128);not null;uniqueIndex:uk_anonymous_id" json:"anonymous_id"`
	Username      *string   `gorm:"column:username;type:varchar(64);uniqueIndex:uk_user_username" json:"username,omitempty"`
	PasswordHash  *string   `gorm:"column:password_hash;type:varchar(255)" json:"-"`
	Nickname      *string   `gorm:"column:nickname;type:varchar(64)" json:"nickname,omitempty"`
	AvatarURL     *string   `gorm:"column:avatar_url;type:varchar(512)" json:"avatar_url,omitempty"`
	Age           *int      `gorm:"column:age" json:"age,omitempty"`
	HomeRegion    *string   `gorm:"column:home_region;type:varchar(64)" json:"home_region,omitempty"`
	CurrentCityID *int64    `gorm:"column:current_city_id;index:idx_user_current_city" json:"current_city_id"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string { return "users" }
