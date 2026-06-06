package repository

import (
	"context"
	"errors"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AITaskRepo struct {
	DB *gorm.DB
}

func NewAITaskRepo(db *gorm.DB) *AITaskRepo {
	return &AITaskRepo{DB: db}
}

func (r *AITaskRepo) Create(ctx context.Context, task *model.AITask) error {
	return r.DB.WithContext(ctx).Create(task).Error
}

func (r *AITaskRepo) FindByID(ctx context.Context, id int64) (*model.AITask, error) {
	var task model.AITask
	err := r.DB.WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *AITaskRepo) FindByIDForUser(ctx context.Context, id, userID int64) (*model.AITask, error) {
	var task model.AITask
	err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *AITaskRepo) ClaimNext(ctx context.Context, taskType string, maxAttempts int) (*model.AITask, error) {
	var claimed *model.AITask
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task model.AITask
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("type = ? AND status IN ? AND attempts < ?", taskType, []string{"queued", "retryable"}, maxAttempts).
			Order("updated_at ASC").
			First(&task).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		task.Status = "running"
		task.Attempts++
		task.Error = nil
		if err := tx.Save(&task).Error; err != nil {
			return err
		}
		claimed = &task
		return nil
	})
	return claimed, err
}

func (r *AITaskRepo) MarkSucceeded(ctx context.Context, id int64, resultURL string) error {
	return r.DB.WithContext(ctx).Model(&model.AITask{}).Where("id = ?", id).Updates(map[string]any{
		"status":     "succeeded",
		"result_url": resultURL,
		"error":      nil,
	}).Error
}

func (r *AITaskRepo) MarkFailed(ctx context.Context, id int64, message string) error {
	return r.DB.WithContext(ctx).Model(&model.AITask{}).Where("id = ?", id).Updates(map[string]any{
		"status": "failed",
		"error":  message,
	}).Error
}

func (r *AITaskRepo) MarkRetryable(ctx context.Context, id int64, message string) error {
	return r.DB.WithContext(ctx).Model(&model.AITask{}).Where("id = ?", id).Updates(map[string]any{
		"status": "retryable",
		"error":  message,
	}).Error
}

func (r *AITaskRepo) QueueRetry(ctx context.Context, id, userID int64) error {
	result := r.DB.WithContext(ctx).Model(&model.AITask{}).
		Where("id = ? AND user_id = ? AND status IN ?", id, userID, []string{"failed", "retryable"}).
		Updates(map[string]any{"status": "queued", "error": nil})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
