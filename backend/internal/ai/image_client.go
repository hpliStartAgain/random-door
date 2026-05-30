package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ImageClient struct {
	baseURL  string
	apiKey   string
	timeout  time.Duration
	maxRetry int
}

func NewImageClient(baseURL, apiKey string, timeout time.Duration) *ImageClient {
	return &ImageClient{
		baseURL:  baseURL,
		apiKey:   apiKey,
		timeout:  timeout,
		maxRetry: 1,
	}
}

type imageRequest struct {
	Prompt         string `json:"prompt"`
	SelfieBase64   string `json:"selfie_base64,omitempty"`
	RefImageBase64 string `json:"ref_image_base64,omitempty"`
}

type imageResponse struct {
	Data []struct {
		B64JSON string `json:"b64_json"`
		URL     string `json:"url"`
	} `json:"data"`
}

// Generate calls the external image generation API.
// Returns raw image bytes on success.
func (c *ImageClient) Generate(ctx context.Context, selfiePath, refImagePath, prompt string) ([]byte, error) {
	selfieB64, err := fileToBase64(selfiePath)
	if err != nil {
		return nil, fmt.Errorf("read selfie: %w", err)
	}

	refB64, err := fileToBase64(refImagePath)
	if err != nil {
		slog.Warn("ref image not found, proceeding without", "path", refImagePath, "error", err)
		refB64 = ""
	}

	reqBody := imageRequest{
		Prompt:         prompt,
		SelfieBase64:   selfieB64,
		RefImageBase64: refB64,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetry; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		imgBytes, err := c.doRequest(ctx, data)
		if err != nil {
			lastErr = err
			if ctx.Err() != nil {
				return nil, ErrAITimeout
			}
			slog.Warn("image generation failed, retrying", "attempt", attempt+1, "error", err)
			continue
		}
		return imgBytes, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrAIUpstream, lastErr)
}

func (c *ImageClient) doRequest(ctx context.Context, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	url := c.baseURL + "/images/generations"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ErrAITimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrAIUpstream, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status=%d body=%s", ErrAIUpstream, resp.StatusCode, string(body))
	}

	var imgResp imageResponse
	if err := json.Unmarshal(body, &imgResp); err != nil {
		return nil, fmt.Errorf("%w: parse response: %v", ErrAIUpstream, err)
	}

	if len(imgResp.Data) == 0 {
		return nil, fmt.Errorf("%w: empty data", ErrAIUpstream)
	}

	// Prefer b64 data, fall back to URL
	if imgResp.Data[0].B64JSON != "" {
		imgBytes, err := base64.StdEncoding.DecodeString(imgResp.Data[0].B64JSON)
		if err != nil {
			return nil, fmt.Errorf("decode b64: %w", err)
		}
		return imgBytes, nil
	}

	// Fetch from URL if b64 not available
	if imgResp.Data[0].URL != "" {
		return fetchImage(ctx, imgResp.Data[0].URL)
	}

	return nil, fmt.Errorf("%w: no image data", ErrAIUpstream)
}

func fileToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func fetchImage(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
