package model

import "time"

type City struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name               string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex:uk_city_name" json:"name"`
	Province           string    `gorm:"column:province;type:varchar(64);not null" json:"province"`
	Lat                float64   `gorm:"column:lat;not null" json:"lat"`
	Lng                float64   `gorm:"column:lng;not null" json:"lng"`
	Intro              *string   `gorm:"column:intro;type:text" json:"intro,omitempty"`
	CoverImageURL      *string   `gorm:"column:cover_image_url;type:varchar(512)" json:"cover_image_url,omitempty"`
	DialectSample      *string   `gorm:"column:dialect_sample;type:varchar(255)" json:"dialect_sample,omitempty"`
	DialectExplanation *string   `gorm:"column:dialect_explanation;type:text" json:"dialect_explanation,omitempty"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (City) TableName() string { return "cities" }

// CityCounts holds aggregated content counts for one city (used by list API).
type CityCounts struct {
	LandmarkCount  int
	FoodCount      int
	CharacterCount int
}
