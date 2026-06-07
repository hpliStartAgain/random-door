package model

import "time"

type GuessChallenge struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code       string    `gorm:"column:code;type:varchar(16);not null;uniqueIndex:uk_guess_challenges_code" json:"code"`
	UserID     *int64    `gorm:"column:user_id;index:idx_guess_challenges_user_time,priority:1" json:"user_id,omitempty"`
	CityID     int64     `gorm:"column:city_id;not null;index:idx_guess_challenges_city" json:"city_id"`
	TargetName *string   `gorm:"column:target_name;type:varchar(128)" json:"target_name,omitempty"`
	ImageURL   *string   `gorm:"column:image_url;type:varchar(512)" json:"image_url,omitempty"`
	Caption    *string   `gorm:"column:caption;type:varchar(300)" json:"caption,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_guess_challenges_user_time,priority:2" json:"created_at"`
	ExpiresAt  time.Time `gorm:"column:expires_at;not null;index:idx_guess_challenges_expires" json:"expires_at"`
}

func (GuessChallenge) TableName() string { return "guess_challenges" }
