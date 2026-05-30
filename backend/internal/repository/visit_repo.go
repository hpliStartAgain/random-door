package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type VisitRepo struct {
	DB *gorm.DB
}

func NewVisitRepo(db *gorm.DB) *VisitRepo {
	return &VisitRepo{DB: db}
}

func (r *VisitRepo) Create(ctx context.Context, visit *model.CityVisit) error {
	return r.DB.WithContext(ctx).Create(visit).Error
}

func (r *VisitRepo) CreateTx(tx *gorm.DB, visit *model.CityVisit) error {
	return tx.Create(visit).Error
}

func (r *VisitRepo) ListVisitedCityIDs(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	err := r.DB.WithContext(ctx).Model(&model.CityVisit{}).
		Where("user_id = ?", userID).
		Distinct("city_id").
		Pluck("city_id", &ids).Error
	return ids, err
}

func (r *VisitRepo) CountByUserMode(ctx context.Context, userID int64, mode string) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.CityVisit{}).
		Where("user_id = ? AND visit_mode = ?", userID, mode).
		Count(&count).Error
	return count, err
}
