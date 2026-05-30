package model

import "time"

type DiceRoll struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"dice_roll_id"`
	UserID     int64     `gorm:"column:user_id;not null;index:idx_dr_user_time" json:"user_id"`
	FromCityID *int64    `gorm:"column:from_city_id;index:idx_dr_from_city" json:"from_city_id,omitempty"`
	ToCityID   int64     `gorm:"column:to_city_id;not null;index:idx_dr_to_city" json:"to_city_id"`
	Direction  string    `gorm:"column:direction;type:varchar(32);not null" json:"direction"`
	DistanceKm int       `gorm:"column:distance_km;not null" json:"distance_km"`
	TargetLat  *float64  `gorm:"column:target_lat" json:"target_lat,omitempty"`
	TargetLng  *float64  `gorm:"column:target_lng" json:"target_lng,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (DiceRoll) TableName() string { return "dice_rolls" }
