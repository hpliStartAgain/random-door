package model

import "time"

type CityTag struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CityID    int64     `gorm:"column:city_id;not null;index:idx_ct_city;uniqueIndex:uk_ct_city_tag,priority:1" json:"city_id"`
	Tag       string    `gorm:"column:tag;type:varchar(64);not null;index:idx_ct_tag;uniqueIndex:uk_ct_city_tag,priority:2" json:"tag"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CityTag) TableName() string { return "city_tags" }
