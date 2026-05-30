package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type CityRepo struct {
	DB *gorm.DB
}

func NewCityRepo(db *gorm.DB) *CityRepo {
	return &CityRepo{DB: db}
}

func (r *CityRepo) ListAll(ctx context.Context) ([]model.City, error) {
	var cities []model.City
	err := r.DB.WithContext(ctx).Order("id ASC").Find(&cities).Error
	return cities, err
}

func (r *CityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	var city model.City
	err := r.DB.WithContext(ctx).First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *CityRepo) ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error) {
	var tags []model.CityTag
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&tags).Error
	return tags, err
}

func (r *CityRepo) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	var landmarks []model.Landmark
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&landmarks).Error
	return landmarks, err
}

func (r *CityRepo) ListFoods(ctx context.Context, cityID int64) ([]model.Food, error) {
	var foods []model.Food
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&foods).Error
	return foods, err
}

func (r *CityRepo) ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error) {
	var chars []model.Character
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&chars).Error
	return chars, err
}

func (r *CityRepo) FindCharacterByID(ctx context.Context, id int64) (*model.Character, error) {
	var ch model.Character
	err := r.DB.WithContext(ctx).First(&ch, id).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}
