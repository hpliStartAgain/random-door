package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type chatCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
	FindCharacterByID(ctx context.Context, id int64) (*model.Character, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
	ListFoods(ctx context.Context, cityID int64) ([]model.Food, error)
}

type chatMessageRepo interface {
	Create(ctx context.Context, msg *model.ChatMessage) error
	ListByUserCharacter(ctx context.Context, userID, characterID int64, limit int) ([]model.ChatMessage, error)
}

type llmCaller interface {
	Chat(ctx context.Context, systemPrompt, userMessage string) (string, error)
	ChatWithHistory(ctx context.Context, systemPrompt string, history []ai.ChatMessage, userMessage string) (string, error)
}

type ChatService struct {
	cityRepo chatCityRepo
	chatRepo chatMessageRepo
	llm      llmCaller
}

func NewChatService(cityRepo chatCityRepo, chatRepo chatMessageRepo, llm llmCaller) *ChatService {
	return &ChatService{cityRepo: cityRepo, chatRepo: chatRepo, llm: llm}
}

const maxMessageLen = 500

// Chat handles a user message to a city character, returns the AI reply.
func (s *ChatService) Chat(ctx context.Context, userID, cityID, characterID int64, message string) (string, error) {
	if len([]rune(message)) > maxMessageLen {
		return "", invalidParam("message too long (max 500 characters)")
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

	// 2. Load landmarks and foods for prompt context
	landmarks, _ := s.cityRepo.ListLandmarks(ctx, cityID)
	lmNames := make([]string, 0, len(landmarks))
	for _, l := range landmarks {
		lmNames = append(lmNames, l.Name)
	}

	foods, _ := s.cityRepo.ListFoods(ctx, cityID)
	foodNames := make([]string, 0, len(foods))
	for _, f := range foods {
		foodNames = append(foodNames, f.Name)
	}

	// 3. Build system prompt
	dialectStyle := ""
	if ch.DialectStyle != nil {
		dialectStyle = *ch.DialectStyle
	}
	systemPrompt := ai.BuildChatPrompt(ai.ChatContext{
		CityName:      city.Name,
		CharacterName: ch.Name,
		Persona:       ch.Persona,
		Landmarks:     lmNames,
		Foods:         foodNames,
		DialectStyle:  dialectStyle,
	})

	// 4. Load history
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

	// 5. Call LLM
	reply, err := s.llm.ChatWithHistory(ctx, systemPrompt, aiHistory, message)
	if err != nil {
		slog.Error("llm chat failed", "error", err, "user_id", userID, "character_id", characterID)
		if errors.Is(err, ai.ErrAITimeout) {
			return "", ai.ErrAITimeout
		}
		return "", fmt.Errorf("%w: %v", ai.ErrAIUpstream, err)
	}

	// 6. Save user and assistant messages
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
