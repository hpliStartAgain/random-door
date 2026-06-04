package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/achievement"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

// CheckinStore writes a checkin and evaluates achievements in a single transaction.
type CheckinStore struct {
	db          *gorm.DB
	checkinRepo *CheckinRepo
}

func NewCheckinStore(db *gorm.DB, checkinRepo *CheckinRepo) *CheckinStore {
	return &CheckinStore{db: db, checkinRepo: checkinRepo}
}

// CreateAndEvaluate persists the checkin then evaluates achievements; both
// succeed or both roll back (api-contract.md §9 requires this to be atomic).
func (s *CheckinStore) CreateAndEvaluate(ctx context.Context, checkin *model.Checkin) ([]model.Achievement, error) {
	var unlocked []model.Achievement
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.checkinRepo.CreateTx(tx, checkin); err != nil {
			return err
		}
		evaluated, e := achievement.Evaluate(ctx, checkin.UserID, achievement.Repos{DB: tx})
		if e != nil {
			return e
		}
		unlocked = evaluated
		return nil
	})
	if err != nil {
		return nil, err
	}
	return unlocked, nil
}
