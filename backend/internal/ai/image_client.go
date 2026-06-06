package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageClient struct {
	baseURL   string
	apiKey    string
	model     string
	staticDir string
	uploadDir string
	timeout   time.Duration
	maxRetry  int
}

func NewImageClient(baseURL, apiKey, model, staticDir, uploadDir string, timeout time.Duration) *ImageClient {
	if model == "" {
		model = "wan2.7-image-pro"
	}
	return &ImageClient{
		baseURL:   baseURL,
		apiKey:    apiKey,
		model:     model,
		staticDir: staticDir,
		uploadDir: uploadDir,
		timeout:   timeout,
		maxRetry:  1,
	}
}

type openAIImageRequest struct {
	Prompt         string `json:"prompt"`
	SelfieBase64   string `json:"selfie_base64,omitempty"`
	RefImageBase64 string `json:"ref_image_base64,omitempty"`
}

type dashScopeRequest struct {
	Model      string              `json:"model"`
	Input      dashScopeInput      `json:"input"`
	Parameters dashScopeParameters `json:"parameters"`
}

type dashScopeInput struct {
	Messages []dashScopeMessage `json:"messages"`
}

type dashScopeMessage struct {
	Role    string             `json:"role"`
	Content []dashScopeContent `json:"content"`
}

type dashScopeContent struct {
	Text  string `json:"text,omitempty"`
	Image string `json:"image,omitempty"`
}

type dashScopeParameters struct {
	Size         string `json:"size"`
	N            int    `json:"n"`
	Watermark    bool   `json:"watermark"`
	ThinkingMode bool   `json:"thinking_mode"`
}

// Generate calls the external image generation API.
// Returns raw image bytes on success.
func (c *ImageClient) Generate(ctx context.Context, selfiePath, refImagePath, prompt string) ([]byte, error) {
	if c.isMock() {
		slog.Info("image generation using mock renderer", "selfie_path", selfiePath, "ref_image_path", refImagePath)
		return mockGeneratedImage()
	}

	selfieDataURL, err := c.fileToDataURL(selfiePath)
	if err != nil {
		return nil, fmt.Errorf("read selfie: %w", err)
	}

	refDataURL, err := c.fileToDataURL(refImagePath)
	if err != nil {
		slog.Warn("ref image not found, proceeding without", "path", refImagePath, "error", err)
		refDataURL = ""
	}

	var reqBody any
	if c.isDashScope() {
		reqBody = c.dashScopePayload(prompt, selfieDataURL, refDataURL)
	} else {
		reqBody = openAIImageRequest{
			Prompt:         prompt,
			SelfieBase64:   stripDataURL(selfieDataURL),
			RefImageBase64: stripDataURL(refDataURL),
		}
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
			if errors.Is(err, ErrAITimeout) {
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

	url := c.endpoint()
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

	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("%w: parse response: %v", ErrAIUpstream, err)
	}

	imgBytes, err := decodeFirstImage(ctx, decoded)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAIUpstream, err)
	}
	return imgBytes, nil
}

func (c *ImageClient) isMock() bool {
	return c == nil ||
		isPlaceholderValue(c.apiKey) ||
		isPlaceholderValue(c.baseURL) ||
		strings.EqualFold(strings.TrimSpace(c.baseURL), "mock")
}

func (c *ImageClient) isDashScope() bool {
	base := strings.ToLower(c.baseURL)
	return strings.Contains(base, "dashscope.aliyuncs.com") || strings.Contains(base, "multimodal-generation")
}

func (c *ImageClient) endpoint() string {
	base := strings.TrimRight(c.baseURL, "/")
	if c.isDashScope() {
		if strings.HasSuffix(base, "/generation") {
			return base
		}
		return base + "/api/v1/services/aigc/multimodal-generation/generation"
	}
	return base + "/images/generations"
}

func (c *ImageClient) dashScopePayload(prompt, selfieDataURL, refDataURL string) dashScopeRequest {
	content := []dashScopeContent{{Text: prompt}}
	if selfieDataURL != "" {
		content = append(content, dashScopeContent{Image: selfieDataURL})
	}
	if refDataURL != "" {
		content = append(content, dashScopeContent{Image: refDataURL})
	}
	return dashScopeRequest{
		Model: c.model,
		Input: dashScopeInput{Messages: []dashScopeMessage{{
			Role:    "user",
			Content: content,
		}}},
		Parameters: dashScopeParameters{
			Size:         "2K",
			N:            1,
			Watermark:    false,
			ThinkingMode: true,
		},
	}
}

func mockGeneratedImage() ([]byte, error) {
	const width = 768
	const height = 1024
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ratioX := float64(x) / width
			ratioY := float64(y) / height
			r := uint8(34 + 42*ratioX)
			g := uint8(48 + 80*ratioY)
			b := uint8(60 + 70*(1-ratioY))
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	gold := color.RGBA{R: 210, G: 170, B: 92, A: 255}
	soft := color.RGBA{R: 246, G: 238, B: 213, A: 230}
	for i := 0; i < 18; i++ {
		drawRect(img, 54+i, 54+i, width-54-i, 58+i, gold)
		drawRect(img, 54+i, height-58-i, width-54-i, height-54-i, gold)
		drawRect(img, 54+i, 54+i, 58+i, height-54-i, gold)
		drawRect(img, width-58-i, 54+i, width-54-i, height-54-i, gold)
	}
	drawCircle(img, width/2, height/2-80, 150, soft)
	drawCircle(img, width/2, height/2-80, 88, color.RGBA{R: 226, G: 198, B: 132, A: 255})
	drawRect(img, 190, 650, 578, 850, color.RGBA{R: 246, G: 238, B: 213, A: 210})
	drawRect(img, 230, 692, 538, 810, color.RGBA{R: 43, G: 58, B: 54, A: 210})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode mock image: %w", err)
	}
	return buf.Bytes(), nil
}

func drawRect(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	bounds := img.Bounds()
	if x1 < bounds.Min.X {
		x1 = bounds.Min.X
	}
	if y1 < bounds.Min.Y {
		y1 = bounds.Min.Y
	}
	if x2 > bounds.Max.X {
		x2 = bounds.Max.X
	}
	if y2 > bounds.Max.Y {
		y2 = bounds.Max.Y
	}
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, radius int, c color.RGBA) {
	r2 := radius * radius
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r2 && image.Pt(x, y).In(img.Bounds()) {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func (c *ImageClient) fileToDataURL(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path is empty")
	}
	localPath := c.resolveLocalPath(path)
	data, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}
	mimeType := mimeTypeFromPath(localPath)
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

func (c *ImageClient) resolveLocalPath(path string) string {
	normalized := filepath.ToSlash(path)
	switch {
	case strings.HasPrefix(normalized, "/uploads/") && c.uploadDir != "":
		return filepath.Join(c.uploadDir, filepath.FromSlash(strings.TrimPrefix(normalized, "/uploads/")))
	case strings.HasPrefix(normalized, "/static/") && c.staticDir != "":
		return filepath.Join(c.staticDir, filepath.FromSlash(strings.TrimPrefix(normalized, "/static/")))
	default:
		return path
	}
}

func mimeTypeFromPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	default:
		return "image/png"
	}
}

func stripDataURL(dataURL string) string {
	if i := strings.Index(dataURL, ","); i >= 0 {
		return dataURL[i+1:]
	}
	return dataURL
}

func decodeFirstImage(ctx context.Context, value any) ([]byte, error) {
	switch v := value.(type) {
	case map[string]any:
		for _, key := range []string{"b64_json", "image", "image_url", "url"} {
			if raw, ok := v[key].(string); ok && raw != "" {
				if data, err := decodeImageValue(ctx, key, raw); err == nil {
					return data, nil
				}
			}
		}
		for _, child := range v {
			if data, err := decodeFirstImage(ctx, child); err == nil {
				return data, nil
			}
		}
	case []any:
		for _, child := range v {
			if data, err := decodeFirstImage(ctx, child); err == nil {
				return data, nil
			}
		}
	}
	return nil, fmt.Errorf("no image data in response")
}

func decodeImageValue(ctx context.Context, key, value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("empty image value")
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return fetchImage(ctx, value)
	}
	if strings.HasPrefix(value, "data:image/") {
		return decodeBase64Image(stripDataURL(value))
	}
	if key == "b64_json" {
		return decodeBase64Image(value)
	}
	return nil, fmt.Errorf("unsupported image value")
}

func decodeBase64Image(value string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err == nil {
		return data, nil
	}
	data, rawErr := base64.RawStdEncoding.DecodeString(value)
	if rawErr == nil {
		return data, nil
	}
	return nil, err
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch image status=%d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
