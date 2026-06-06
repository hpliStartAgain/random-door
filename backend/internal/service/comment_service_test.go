package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeCommentRepo struct {
	create       func(context.Context, *model.Comment) error
	listByTarget func(context.Context, string, int64, int) ([]model.Comment, error)
}

func (r *fakeCommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	if r.create == nil {
		return nil
	}
	return r.create(ctx, comment)
}

func (r *fakeCommentRepo) ListByTarget(ctx context.Context, targetType string, targetID int64, limit int) ([]model.Comment, error) {
	if r.listByTarget == nil {
		return nil, nil
	}
	return r.listByTarget(ctx, targetType, targetID, limit)
}

type fakeCommentTargetRepo struct {
	findLandmarkByID  func(context.Context, int64) (*model.Landmark, error)
	findFoodByID      func(context.Context, int64) (*model.Food, error)
	findCharacterByID func(context.Context, int64) (*model.Character, error)
}

func (r *fakeCommentTargetRepo) FindLandmarkByID(ctx context.Context, id int64) (*model.Landmark, error) {
	if r.findLandmarkByID == nil {
		return &model.Landmark{ID: id}, nil
	}
	return r.findLandmarkByID(ctx, id)
}

func (r *fakeCommentTargetRepo) FindFoodByID(ctx context.Context, id int64) (*model.Food, error) {
	if r.findFoodByID == nil {
		return &model.Food{ID: id}, nil
	}
	return r.findFoodByID(ctx, id)
}

func (r *fakeCommentTargetRepo) FindCharacterByID(ctx context.Context, id int64) (*model.Character, error) {
	if r.findCharacterByID == nil {
		return &model.Character{ID: id}, nil
	}
	return r.findCharacterByID(ctx, id)
}

func TestCommentServiceCreateDefaultsNicknameAndTrimsContent(t *testing.T) {
	var saved *model.Comment
	svc := NewCommentService(&fakeCommentRepo{
		create: func(_ context.Context, comment *model.Comment) error {
			comment.ID = 7
			comment.CreatedAt = time.Date(2026, 6, 6, 9, 0, 0, 0, time.UTC)
			saved = comment
			return nil
		},
	}, &fakeCommentTargetRepo{})

	got, err := svc.Create(context.Background(), CreateCommentRequest{
		TargetType: CommentTargetLandmark,
		TargetID:   12,
		Content:    "  这个角度适合打卡  ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if saved == nil || saved.Nickname != "游客" || saved.Content != "这个角度适合打卡" {
		t.Fatalf("saved comment = %#v, want trimmed default nickname comment", saved)
	}
	if got.ID != 7 || got.Nickname != "游客" || got.CreatedAt == "" {
		t.Fatalf("Create() = %#v, want public comment item", got)
	}
}

func TestCommentServiceListReturnsOldestFirstForDanmaku(t *testing.T) {
	svc := NewCommentService(&fakeCommentRepo{
		listByTarget: func(_ context.Context, targetType string, targetID int64, limit int) ([]model.Comment, error) {
			if targetType != CommentTargetFood || targetID != 5 || limit != defaultCommentLimit {
				t.Fatalf("ListByTarget args = %s %d %d", targetType, targetID, limit)
			}
			return []model.Comment{
				{ID: 2, TargetType: targetType, TargetID: targetID, Nickname: "乙", Content: "后发"},
				{ID: 1, TargetType: targetType, TargetID: targetID, Nickname: "甲", Content: "先发"},
			}, nil
		},
	}, &fakeCommentTargetRepo{})

	got, err := svc.List(context.Background(), CommentTargetFood, 5, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.Comments) != 2 || got.Comments[0].ID != 1 || got.Comments[1].ID != 2 {
		t.Fatalf("List() = %#v, want oldest-first order", got.Comments)
	}
}

func TestCommentServiceRejectsInvalidInputs(t *testing.T) {
	svc := NewCommentService(&fakeCommentRepo{}, &fakeCommentTargetRepo{})

	tests := []CreateCommentRequest{
		{TargetType: "city", TargetID: 1, Content: "hi"},
		{TargetType: CommentTargetCharacter, TargetID: 0, Content: "hi"},
		{TargetType: CommentTargetCharacter, TargetID: 1, UserID: int64Ptr(-1), Content: "hi"},
		{TargetType: CommentTargetCharacter, TargetID: 1, Content: ""},
		{TargetType: CommentTargetCharacter, TargetID: 1, Content: strings.Repeat("好", 201)},
	}
	for _, tc := range tests {
		_, err := svc.Create(context.Background(), tc)
		if !errors.Is(err, ErrInvalidParam) {
			t.Fatalf("Create(%#v) error = %v, want ErrInvalidParam", tc, err)
		}
	}
}

func TestCommentServiceMissingTarget(t *testing.T) {
	svc := NewCommentService(&fakeCommentRepo{}, &fakeCommentTargetRepo{
		findCharacterByID: func(context.Context, int64) (*model.Character, error) {
			return nil, gorm.ErrRecordNotFound
		},
	})

	_, err := svc.Create(context.Background(), CreateCommentRequest{
		TargetType: CommentTargetCharacter,
		TargetID:   8,
		Content:    "想听你讲讲这座城",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Create() error = %v, want ErrNotFound", err)
	}
}

func int64Ptr(v int64) *int64 {
	return &v
}
