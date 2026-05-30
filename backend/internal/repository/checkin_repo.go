package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type CheckinRepo struct {
	DB *gorm.DB
}

func NewCheckinRepo(db *gorm.DB) *CheckinRepo {
	return &CheckinRepo{DB: db}
}

func (r *CheckinRepo) Create(ctx context.Context, checkin *model.Checkin) error {
	return r.DB.WithContext(ctx).Create(checkin).Error
}

func (r *CheckinRepo) CreateTx(tx *gorm.DB, checkin *model.Checkin) error {
	return tx.Create(checkin).Error
}

func (r *CheckinRepo) CountByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.Checkin{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *CheckinRepo) ListByUser(ctx context.Context, userID int64) ([]model.Checkin, error) {
	var list []model.Checkin
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Find(&list).Error
	return list, err
}
