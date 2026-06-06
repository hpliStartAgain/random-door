package repository

import (
	"context"
	"errors"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AIUsageRepo struct {
	DB *gorm.DB
}

func NewAIUsageRepo(db *gorm.DB) *AIUsageRepo {
	return &AIUsageRepo{DB: db}
}

func (r *AIUsageRepo) IncrementIfBelow(ctx context.Context, userID int64, usageType string, usageDate time.Time, limit int) (int, bool, error) {
	if limit <= 0 {
		return 0, true, nil
	}

	var nextCount int
	allowed := false
	day := truncateDate(usageDate)
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var log model.AIUsageLog
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND usage_type = ? AND usage_date = ?", userID, usageType, day).
			First(&log).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log = model.AIUsageLog{UserID: userID, UsageType: usageType, UsageDate: day, Count: 1}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
			nextCount = 1
			allowed = true
			return nil
		}
		if err != nil {
			return err
		}
		if log.Count >= limit {
			nextCount = log.Count
			allowed = false
			return nil
		}
		log.Count++
		if err := tx.Save(&log).Error; err != nil {
			return err
		}
		nextCount = log.Count
		allowed = true
		return nil
	})
	return nextCount, allowed, err
}

func truncateDate(t time.Time) time.Time {
	year, month, day := t.Local().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Local().Location())
}
