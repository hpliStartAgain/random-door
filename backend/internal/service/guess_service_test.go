package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeGuessCityRepo struct {
	findByID func(context.Context, int64) (*model.City, error)
}

func (r *fakeGuessCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}

func (r *fakeGuessCityRepo) ListLandmarks(context.Context, int64) ([]model.Landmark, error) {
	return []model.Landmark{{Name: "城墙"}}, nil
}

func (r *fakeGuessCityRepo) ListFoods(context.Context, int64) ([]model.Food, error) {
	return []model.Food{{Name: "肉夹馍"}}, nil
}

func (r *fakeGuessCityRepo) ListCharacters(context.Context, int64) ([]model.Character, error) {
	return []model.Character{{Name: "李白"}}, nil
}

type fakeGuessLLM struct {
	replies []string
	err     error
	calls   int
}

func (f *fakeGuessLLM) Chat(context.Context, string, string) (string, error) {
	f.calls++
	if f.err != nil {
		return "", f.err
	}
	if len(f.replies) >= f.calls {
		return f.replies[f.calls-1], nil
	}
	return "生成文案", nil
}

func TestGuessServiceGenerateCaptionLLM(t *testing.T) {
	llm := &fakeGuessLLM{replies: []string{"微博文案 #任意门漫游#", "朋友圈文案"}}
	svc := NewGuessService(okGuessCityRepo(), llm)

	got, err := svc.GenerateCaption(context.Background(), GuessCaptionRequest{CityID: 1, TargetName: "城墙"})
	if err != nil {
		t.Fatalf("GenerateCaption() error = %v", err)
	}
	if got.Weibo != "微博文案 #任意门漫游#" || got.Moments != "朋友圈文案" || llm.calls != 2 {
		t.Fatalf("GenerateCaption() = %#v, calls=%d", got, llm.calls)
	}
}

func TestGuessServiceGenerateCaptionFallback(t *testing.T) {
	svc := NewGuessService(okGuessCityRepo(), &fakeGuessLLM{err: errors.New("llm down")})

	got, err := svc.GenerateCaption(context.Background(), GuessCaptionRequest{CityID: 1})
	if err != nil {
		t.Fatalf("GenerateCaption() error = %v", err)
	}
	if !strings.Contains(got.Weibo, "西安") || !strings.Contains(got.Moments, "西安") {
		t.Fatalf("fallback captions = %#v, want city name", got)
	}
}

func TestGuessServiceCityNotFound(t *testing.T) {
	svc := NewGuessService(&fakeGuessCityRepo{
		findByID: func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
	}, nil)

	_, err := svc.GenerateCaption(context.Background(), GuessCaptionRequest{CityID: 404})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GenerateCaption() error = %v, want ErrNotFound", err)
	}
}

func TestGuessServiceSceneHintTooLong(t *testing.T) {
	svc := NewGuessService(okGuessCityRepo(), nil)
	_, err := svc.GenerateCaption(context.Background(), GuessCaptionRequest{
		CityID:    1,
		SceneHint: strings.Repeat("好", 201),
	})
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("GenerateCaption() error = %v, want ErrInvalidParam", err)
	}
}

func okGuessCityRepo() *fakeGuessCityRepo {
	return &fakeGuessCityRepo{
		findByID: func(context.Context, int64) (*model.City, error) {
			return &model.City{ID: 1, Name: "西安", Province: "陕西"}, nil
		},
	}
}
