package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type AchievementRepo struct {
	DB *gorm.DB
}

func NewAchievementRepo(db *gorm.DB) *AchievementRepo {
	return &AchievementRepo{DB: db}
}

func (r *AchievementRepo) ListAll(ctx context.Context) ([]model.Achievement, error) {
	var list []model.Achievement
	err := r.DB.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *AchievementRepo) ListUserAchievements(ctx context.Context, userID int64) ([]model.UserAchievement, error) {
	var list []model.UserAchievement
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *AchievementRepo) CreateUserAchievement(ctx context.Context, ua *model.UserAchievement) error {
	return r.DB.WithContext(ctx).Create(ua).Error
}
