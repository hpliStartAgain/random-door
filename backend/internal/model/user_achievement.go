package model

import "time"

type UserAchievement struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"column:user_id;not null;uniqueIndex:uk_ua_user_ach" json:"user_id"`
	AchievementID int64     `gorm:"column:achievement_id;not null;uniqueIndex:uk_ua_user_ach;index:idx_ua_achievement" json:"achievement_id"`
	UnlockedAt    time.Time `gorm:"column:unlocked_at;not null" json:"unlocked_at"`
}

func (UserAchievement) TableName() string { return "user_achievements" }
