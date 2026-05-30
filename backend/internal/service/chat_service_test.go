package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeChatCityRepo struct {
	findByID          func(context.Context, int64) (*model.City, error)
	findCharacterByID func(context.Context, int64) (*model.Character, error)
	listLandmarks     func(context.Context, int64) ([]model.Landmark, error)
	listFoods         func(context.Context, int64) ([]model.Food, error)
}

func (r *fakeChatCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}
func (r *fakeChatCityRepo) FindCharacterByID(ctx context.Context, id int64) (*model.Character, error) {
	return r.findCharacterByID(ctx, id)
}
func (r *fakeChatCityRepo) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	if r.listLandmarks == nil {
		return nil, nil
	}
	return r.listLandmarks(ctx, cityID)
}
func (r *fakeChatCityRepo) ListFoods(ctx context.Context, cityID int64) ([]model.Food, error) {
	if r.listFoods == nil {
		return nil, nil
	}
	return r.listFoods(ctx, cityID)
}

type fakeChatMsgRepo struct {
	create              func(context.Context, *model.ChatMessage) error
	listByUserCharacter func(context.Context, int64, int64, int) ([]model.ChatMessage, error)
}

func (r *fakeChatMsgRepo) Create(ctx context.Context, msg *model.ChatMessage) error {
	if r.create == nil {
		return nil
	}
	return r.create(ctx, msg)
}

func (r *fakeChatMsgRepo) ListByUserCharacter(ctx context.Context, userID, characterID int64, limit int) ([]model.ChatMessage, error) {
	if r.listByUserCharacter == nil {
		return nil, nil
	}
	return r.listByUserCharacter(ctx, userID, characterID, limit)
}

type fakeLLM struct {
	chat            func(context.Context, string, string) (string, error)
	chatWithHistory func(context.Context, string, []aiPkg.ChatMessage, string) (string, error)
}

func (f *fakeLLM) Chat(ctx context.Context, sys, user string) (string, error) {
	return f.chat(ctx, sys, user)
}

func (f *fakeLLM) ChatWithHistory(ctx context.Context, sys string, hist []aiPkg.ChatMessage, user string) (string, error) {
	if f.chatWithHistory != nil {
		return f.chatWithHistory(ctx, sys, hist, user)
	}
	if f.chat != nil {
		return f.chat(ctx, sys, user)
	}
	return "", nil
}

func newTestChatService(cityRepo chatCityRepo, chatRepo chatMessageRepo, llm llmCaller) *ChatService {
	return &ChatService{cityRepo: cityRepo, chatRepo: chatRepo, llm: llm}
}

var testCity = &model.City{ID: 1, Name: "西安", Province: "陕西"}
var testChar = &model.Character{ID: 8, CityID: 1, Name: "李白", Persona: "诗仙"}

func okCityRepo() *fakeChatCityRepo {
	return &fakeChatCityRepo{
		findByID:          func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
		findCharacterByID: func(_ context.Context, _ int64) (*model.Character, error) { return testChar, nil },
	}
}

func TestChatService_Success(t *testing.T) {
	svc := newTestChatService(okCityRepo(), &fakeChatMsgRepo{}, &fakeLLM{
		chat: func(_ context.Context, _, _ string) (string, error) { return "长安风物", nil },
	})
	reply, err := svc.Chat(context.Background(), 1, 1, 8, "你好")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "长安风物" {
		t.Fatalf("unexpected reply: %s", reply)
	}
}

func TestChatService_MessageTooLong(t *testing.T) {
	svc := newTestChatService(okCityRepo(), &fakeChatMsgRepo{}, &fakeLLM{})
	_, err := svc.Chat(context.Background(), 1, 1, 8, strings.Repeat("a", 501))
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("expected ErrInvalidParam, got %v", err)
	}
}

func TestChatService_CityNotFound(t *testing.T) {
	repo := &fakeChatCityRepo{
		findByID:          func(_ context.Context, _ int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
		findCharacterByID: func(_ context.Context, _ int64) (*model.Character, error) { return testChar, nil },
	}
	svc := newTestChatService(repo, &fakeChatMsgRepo{}, &fakeLLM{})
	_, err := svc.Chat(context.Background(), 1, 99, 8, "你好")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestChatService_CharacterNotFound(t *testing.T) {
	repo := &fakeChatCityRepo{
		findByID:          func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
		findCharacterByID: func(_ context.Context, _ int64) (*model.Character, error) { return nil, gorm.ErrRecordNotFound },
	}
	svc := newTestChatService(repo, &fakeChatMsgRepo{}, &fakeLLM{})
	_, err := svc.Chat(context.Background(), 1, 1, 99, "你好")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestChatService_LLMTimeout(t *testing.T) {
	svc := newTestChatService(okCityRepo(), &fakeChatMsgRepo{}, &fakeLLM{
		chat: func(_ context.Context, _, _ string) (string, error) { return "", aiPkg.ErrAITimeout },
	})
	_, err := svc.Chat(context.Background(), 1, 1, 8, "你好")
	if !errors.Is(err, aiPkg.ErrAITimeout) {
		t.Fatalf("expected ErrAITimeout, got %v", err)
	}
}

func TestChatService_LLMUpstreamError(t *testing.T) {
	svc := newTestChatService(okCityRepo(), &fakeChatMsgRepo{}, &fakeLLM{
		chat: func(_ context.Context, _, _ string) (string, error) { return "", aiPkg.ErrAIUpstream },
	})
	_, err := svc.Chat(context.Background(), 1, 1, 8, "你好")
	if !errors.Is(err, aiPkg.ErrAIUpstream) {
		t.Fatalf("expected ErrAIUpstream, got %v", err)
	}
}
