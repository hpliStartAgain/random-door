package ai

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLLMClientPlaceholderConfigFailsFast(t *testing.T) {
	client := NewLLMClient("https://api.your-llm-provider.com/v1", "sk-REPLACE_ME", "deepseek-v4-flash", time.Second)

	_, err := client.Chat(context.Background(), "system", "hello")
	if !errors.Is(err, ErrAIUpstream) {
		t.Fatalf("expected ErrAIUpstream, got %v", err)
	}
}
