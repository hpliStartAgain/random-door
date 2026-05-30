package model

import "time"

type Checkin struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement" json:"checkin_id"`
	UserID            int64     `gorm:"column:user_id;not null;index:idx_ck_user_time;index:idx_ck_user_city" json:"user_id"`
	CityID            int64     `gorm:"column:city_id;not null;index:idx_ck_user_city;index:idx_ck_city" json:"city_id"`
	LandmarkID        *int64    `gorm:"column:landmark_id;index:idx_ck_landmark" json:"landmark_id,omitempty"`
	VisitID           *int64    `gorm:"column:visit_id;index:idx_ck_visit" json:"visit_id,omitempty"`
	GeneratedImageURL *string   `gorm:"column:generated_image_url;type:varchar(512)" json:"generated_image_url,omitempty"`
	CheckinMode       *string   `gorm:"column:checkin_mode;type:varchar(32)" json:"checkin_mode,omitempty"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Checkin) TableName() string { return "checkins" }
