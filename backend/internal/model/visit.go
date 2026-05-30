package model

import "time"

type CityVisit struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"visit_id"`
	UserID     int64     `gorm:"column:user_id;not null;index:idx_cv_user_city;index:idx_cv_user_time" json:"user_id"`
	CityID     int64     `gorm:"column:city_id;not null;index:idx_cv_user_city;index:idx_cv_city" json:"city_id"`
	VisitMode  string    `gorm:"column:visit_mode;type:varchar(32);not null" json:"visit_mode"`
	Source     *string   `gorm:"column:source;type:varchar(64)" json:"source,omitempty"`
	FromCityID *int64    `gorm:"column:from_city_id;index:idx_cv_from_city" json:"from_city_id,omitempty"`
	DiceRollID *int64    `gorm:"column:dice_roll_id;index:idx_cv_dice_roll" json:"dice_roll_id,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CityVisit) TableName() string { return "city_visits" }
