package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{DB: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	return r.DB.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepo) ListByTarget(ctx context.Context, targetType string, targetID int64, limit int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.DB.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("created_at DESC, id DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}
