package model

import "time"

type Food struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CityID      int64     `gorm:"column:city_id;not null;index:idx_food_city;uniqueIndex:uk_food_city_name,priority:1" json:"city_id"`
	Name        string    `gorm:"column:name;type:varchar(128);not null;uniqueIndex:uk_food_city_name,priority:2" json:"name"`
	ImageURL    *string   `gorm:"column:image_url;type:varchar(512)" json:"image_url,omitempty"`
	Description *string   `gorm:"column:description;type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Food) TableName() string { return "foods" }
