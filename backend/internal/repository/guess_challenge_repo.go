package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type GuessChallengeRepo struct {
	DB *gorm.DB
}

func NewGuessChallengeRepo(db *gorm.DB) *GuessChallengeRepo {
	return &GuessChallengeRepo{DB: db}
}

func (r *GuessChallengeRepo) CreateChallenge(ctx context.Context, challenge *model.GuessChallenge) error {
	return r.DB.WithContext(ctx).Create(challenge).Error
}

func (r *GuessChallengeRepo) FindChallengeByCode(ctx context.Context, code string) (*model.GuessChallenge, error) {
	var challenge model.GuessChallenge
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&challenge).Error
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (r *GuessChallengeRepo) CreateAnswer(ctx context.Context, answer *model.GuessAnswer) error {
	return r.DB.WithContext(ctx).Create(answer).Error
}
