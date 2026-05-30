package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

// GameStore persists a dice roll and its resulting visit in a single transaction.
type GameStore struct {
	db        *gorm.DB
	diceRepo  *DiceRepo
	visitRepo *VisitRepo
}

func NewGameStore(db *gorm.DB, diceRepo *DiceRepo, visitRepo *VisitRepo) *GameStore {
	return &GameStore{db: db, diceRepo: diceRepo, visitRepo: visitRepo}
}

// CreateRollWithVisit writes the dice roll first, links the visit to the new
// roll id, then writes the visit. Both succeed or both roll back.
func (s *GameStore) CreateRollWithVisit(ctx context.Context, roll *model.DiceRoll, visit *model.CityVisit) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.diceRepo.CreateTx(tx, roll); err != nil {
			return err
		}
		visit.DiceRollID = &roll.ID
		return s.visitRepo.CreateTx(tx, visit)
	})
}
