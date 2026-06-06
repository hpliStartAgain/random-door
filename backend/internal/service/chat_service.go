package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type chatCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
	FindCharacterByID(ctx context.Context, id int64) (*model.Character, error)
}

type chatMessageRepo interface {
	Create(ctx context.Context, msg *model.ChatMessage) error
	ListByUserCharacter(ctx context.Context, userID, characterID int64, limit int) ([]model.ChatMessage, error)
}

type llmCaller interface {
	Chat(ctx context.Context, systemPrompt, userMessage string) (string, error)
	ChatWithHistory(ctx context.Context, systemPrompt string, history []ai.ChatMessage, userMessage string) (string, error)
}

type aiUsageLimiter interface {
	IncrementIfBelow(ctx context.Context, userID int64, usageType string, usageDate time.Time, limit int) (int, bool, error)
}

type ChatService struct {
	cityRepo       chatCityRepo
	chatRepo       chatMessageRepo
	llm            llmCaller
	usageLimiter   aiUsageLimiter
	chatDailyLimit int
}

func NewChatService(cityRepo chatCityRepo, chatRepo chatMessageRepo, llm llmCaller) *ChatService {
	return &ChatService{cityRepo: cityRepo, chatRepo: chatRepo, llm: llm}
}

func (s *ChatService) WithUsageLimit(limiter aiUsageLimiter, dailyLimit int) *ChatService {
	s.usageLimiter = limiter
	s.chatDailyLimit = dailyLimit
	return s
}

const maxMessageLen = 500

// Chat handles a user message to a city character, returns the AI reply.
func (s *ChatService) Chat(ctx context.Context, userID, cityID, characterID int64, message string) (string, error) {
	if len([]rune(message)) > maxMessageLen {
		return "", invalidParam("message too long (max 500 characters)")
	}
	if s.usageLimiter != nil && s.chatDailyLimit > 0 {
		_, allowed, err := s.usageLimiter.IncrementIfBelow(ctx, userID, "chat", time.Now(), s.chatDailyLimit)
		if err != nil {
			return "", fmt.Errorf("increment chat usage: %w", err)
		}
		if !allowed {
			return "", quotaExceeded("daily chat quota exceeded")
		}
	}
	// 1. Load city and character
	city, err := s.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", notFound("city not found")
		}
		return "", fmt.Errorf("find city: %w", err)
	}

	ch, err := s.cityRepo.FindCharacterByID(ctx, characterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", notFound("character not found")
		}
		return "", fmt.Errorf("find character: %w", err)
	}
	if ch.CityID != cityID {
		return "", notFound("character not found")
	}

	// 2. Build system prompt
	dialectStyle := ""
	if ch.DialectStyle != nil {
		dialectStyle = *ch.DialectStyle
	}
	systemPrompt := ai.BuildChatPrompt(ai.ChatContext{
		CityName:      city.Name,
		CharacterName: ch.Name,
		DialectStyle:  dialectStyle,
	})

	// 3. Load history
	histMsgs, err := s.chatRepo.ListByUserCharacter(ctx, userID, characterID, 10)
	if err != nil {
		slog.Warn("failed to load chat history", "error", err)
	}

	// Reverse history because it's loaded DESC
	var aiHistory []ai.ChatMessage
	for i := len(histMsgs) - 1; i >= 0; i-- {
		aiHistory = append(aiHistory, ai.ChatMessage{
			Role:    histMsgs[i].Role,
			Content: histMsgs[i].Content,
		})
	}

	// 4. Call LLM
	reply, err := s.llm.ChatWithHistory(ctx, systemPrompt, aiHistory, message)
	if err != nil {
		slog.Error("llm chat failed", "error", err, "user_id", userID, "character_id", characterID)
		if errors.Is(err, ai.ErrAITimeout) {
			return "", ai.ErrAITimeout
		}
		return "", fmt.Errorf("%w: %v", ai.ErrAIUpstream, err)
	}

	// 5. Save user and assistant messages
	userMsg := &model.ChatMessage{
		UserID: userID, CityID: cityID, CharacterID: characterID,
		Role: "user", Content: message,
	}
	if err := s.chatRepo.Create(ctx, userMsg); err != nil {
		slog.Warn("save user message failed", "error", err)
	}

	assistantMsg := &model.ChatMessage{
		UserID: userID, CityID: cityID, CharacterID: characterID,
		Role: "assistant", Content: reply,
	}
	if err := s.chatRepo.Create(ctx, assistantMsg); err != nil {
		slog.Warn("save assistant message failed", "error", err)
	}

	slog.Info("chat completed", "user_id", userID, "character", ch.Name, "reply_len", len(reply))
	return reply, nil
}
