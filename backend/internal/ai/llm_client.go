package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrAITimeout  = errors.New("ai: request timeout")
	ErrAIUpstream = errors.New("ai: upstream error")
)

type LLMClient struct {
	baseURL  string
	apiKey   string
	model    string
	timeout  time.Duration
	maxRetry int
}

func NewLLMClient(baseURL, apiKey, model string, timeout time.Duration) *LLMClient {
	return &LLMClient{
		baseURL:  baseURL,
		apiKey:   apiKey,
		model:    model,
		timeout:  timeout,
		maxRetry: 2,
	}
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Chat sends a message to the LLM with the given system prompt and returns the reply.
func (c *LLMClient) Chat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	return c.ChatWithHistory(ctx, systemPrompt, nil, userMessage)
}

// ChatWithHistory sends a message with history to the LLM.
func (c *LLMClient) ChatWithHistory(ctx context.Context, systemPrompt string, history []ChatMessage, userMessage string) (string, error) {
	messages := []ChatMessage{{Role: "system", Content: systemPrompt}}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: userMessage})

	reqBody := chatRequest{
		Model:    c.model,
		Messages: messages,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetry; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			time.Sleep(time.Duration(attempt*500) * time.Millisecond)
		}

		reply, err := c.doRequest(ctx, data)
		if err != nil {
			lastErr = err
			if errors.Is(err, ErrAITimeout) {
				return "", err // Don't retry on timeout
			}
			slog.Warn("llm request failed, retrying", "attempt", attempt+1, "error", err)
			continue
		}
		return reply, nil
	}
	return "", fmt.Errorf("%w: %v", ErrAIUpstream, lastErr)
}

func (c *LLMClient) doRequest(ctx context.Context, data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", ErrAITimeout
		}
		return "", fmt.Errorf("%w: %v", ErrAIUpstream, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status=%d body=%s", ErrAIUpstream, resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("%w: parse response: %v", ErrAIUpstream, err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("%w: empty choices", ErrAIUpstream)
	}

	reply := chatResp.Choices[0].Message.Content
	slog.Info("llm chat completed", "reply_len", len(reply))
	return reply, nil
}
