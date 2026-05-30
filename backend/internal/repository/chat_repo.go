package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type ChatRepo struct {
	DB *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{DB: db}
}

func (r *ChatRepo) Create(ctx context.Context, msg *model.ChatMessage) error {
	return r.DB.WithContext(ctx).Create(msg).Error
}

func (r *ChatRepo) ListByUserCharacter(ctx context.Context, userID, characterID int64, limit int) ([]model.ChatMessage, error) {
	var msgs []model.ChatMessage
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND character_id = ?", userID, characterID).
		Order("created_at DESC").
		Limit(limit).
		Find(&msgs).Error
	return msgs, err
}
