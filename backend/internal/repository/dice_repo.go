package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type DiceRepo struct {
	DB *gorm.DB
}

func NewDiceRepo(db *gorm.DB) *DiceRepo {
	return &DiceRepo{DB: db}
}

func (r *DiceRepo) Create(ctx context.Context, roll *model.DiceRoll) error {
	return r.DB.WithContext(ctx).Create(roll).Error
}

func (r *DiceRepo) CreateTx(tx *gorm.DB, roll *model.DiceRoll) error {
	return tx.Create(roll).Error
}

func (r *DiceRepo) ListRecentByUser(ctx context.Context, userID int64, limit int) ([]model.DiceRoll, error) {
	var rolls []model.DiceRoll
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&rolls).Error
	return rolls, err
}
