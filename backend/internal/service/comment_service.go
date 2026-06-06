package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

const (
	CommentTargetLandmark  = "landmark"
	CommentTargetFood      = "food"
	CommentTargetCharacter = "character"
	defaultCommentLimit    = 50
	maxCommentLimit        = 100
	maxCommentContentLen   = 200
	maxCommentNicknameLen  = 32
)

type commentRepo interface {
	Create(ctx context.Context, comment *model.Comment) error
	ListByTarget(ctx context.Context, targetType string, targetID int64, limit int) ([]model.Comment, error)
}

type commentTargetRepo interface {
	FindLandmarkByID(ctx context.Context, id int64) (*model.Landmark, error)
	FindFoodByID(ctx context.Context, id int64) (*model.Food, error)
	FindCharacterByID(ctx context.Context, id int64) (*model.Character, error)
}

type CommentService struct {
	repo       commentRepo
	targetRepo commentTargetRepo
}

func NewCommentService(repo commentRepo, targetRepo commentTargetRepo) *CommentService {
	return &CommentService{repo: repo, targetRepo: targetRepo}
}

type CreateCommentRequest struct {
	TargetType string `json:"target_type"`
	TargetID   int64  `json:"target_id"`
	UserID     *int64 `json:"user_id,omitempty"`
	Nickname   string `json:"nickname"`
	Content    string `json:"content"`
}

type CommentItem struct {
	ID         int64  `json:"id"`
	TargetType string `json:"target_type"`
	TargetID   int64  `json:"target_id"`
	UserID     *int64 `json:"user_id,omitempty"`
	Nickname   string `json:"nickname"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

type CommentListResult struct {
	Comments []CommentItem `json:"comments"`
}

func (s *CommentService) List(ctx context.Context, targetType string, targetID int64, limit int) (*CommentListResult, error) {
	if err := validateCommentTarget(targetType, targetID); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = defaultCommentLimit
	}
	if limit > maxCommentLimit {
		return nil, invalidParam("limit must be between 1 and 100")
	}
	if err := s.ensureTargetExists(ctx, targetType, targetID); err != nil {
		return nil, err
	}

	rows, err := s.repo.ListByTarget(ctx, targetType, targetID, limit)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	items := make([]CommentItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, commentItem(row))
	}
	return &CommentListResult{Comments: items}, nil
}

func (s *CommentService) Create(ctx context.Context, req CreateCommentRequest) (*CommentItem, error) {
	req.TargetType = strings.TrimSpace(req.TargetType)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Content = strings.TrimSpace(req.Content)

	if err := validateCommentTarget(req.TargetType, req.TargetID); err != nil {
		return nil, err
	}
	if req.UserID != nil && *req.UserID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	if req.Nickname == "" {
		req.Nickname = "游客"
	}
	if len([]rune(req.Nickname)) > maxCommentNicknameLen {
		return nil, invalidParam("nickname too long (max 32 characters)")
	}
	if req.Content == "" {
		return nil, invalidParam("content cannot be empty")
	}
	if len([]rune(req.Content)) > maxCommentContentLen {
		return nil, invalidParam("content too long (max 200 characters)")
	}
	if err := s.ensureTargetExists(ctx, req.TargetType, req.TargetID); err != nil {
		return nil, err
	}

	row := &model.Comment{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		UserID:     req.UserID,
		Nickname:   req.Nickname,
		Content:    req.Content,
	}
	if err := s.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	item := commentItem(*row)
	return &item, nil
}

func validateCommentTarget(targetType string, targetID int64) error {
	switch targetType {
	case CommentTargetLandmark, CommentTargetFood, CommentTargetCharacter:
	default:
		return invalidParam("target_type must be landmark, food, or character")
	}
	if targetID <= 0 {
		return invalidParam("target_id must be a positive integer")
	}
	return nil
}

func (s *CommentService) ensureTargetExists(ctx context.Context, targetType string, targetID int64) error {
	var err error
	switch targetType {
	case CommentTargetLandmark:
		_, err = s.targetRepo.FindLandmarkByID(ctx, targetID)
	case CommentTargetFood:
		_, err = s.targetRepo.FindFoodByID(ctx, targetID)
	case CommentTargetCharacter:
		_, err = s.targetRepo.FindCharacterByID(ctx, targetID)
	}
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("comment target not found")
	}
	return fmt.Errorf("find comment target: %w", err)
}

func commentItem(row model.Comment) CommentItem {
	createdAt := row.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	return CommentItem{
		ID:         row.ID,
		TargetType: row.TargetType,
		TargetID:   row.TargetID,
		UserID:     row.UserID,
		Nickname:   row.Nickname,
		Content:    row.Content,
		CreatedAt:  createdAt.Format(time.RFC3339),
	}
}
