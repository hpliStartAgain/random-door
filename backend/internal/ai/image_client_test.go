package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestImageClientMockGenerateReturnsPNG(t *testing.T) {
	client := NewImageClient("mock", "", "", "", "", time.Second)
	data, err := client.Generate(context.Background(), "", "", "prompt")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("mock image is not a png: %v", err)
	}
}

func TestImageClientPlaceholderConfigUsesMock(t *testing.T) {
	client := NewImageClient("https://api.your-image-provider.com/v1", "sk-REPLACE_ME", "", "", "", time.Second)
	data, err := client.Generate(context.Background(), "", "", "prompt")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("placeholder mock image is not a png: %v", err)
	}
}

func TestImageClientDashScopePayloadAndResponse(t *testing.T) {
	temp := t.TempDir()
	uploadDir := filepath.Join(temp, "uploads")
	if err := os.MkdirAll(filepath.Join(uploadDir, "selfies"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "selfies", "me.png"), tinyPNG(t), 0o644); err != nil {
		t.Fatal(err)
	}

	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/services/aigc/multimodal-generation/generation" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Fatalf("authorization = %q", auth)
		}
		body, _ := ioReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"output":{"choices":[{"message":{"content":[{"image":"data:image/png;base64,` + base64.StdEncoding.EncodeToString(tinyPNG(t)) + `"}]}}]}}`))
	}))
	defer server.Close()

	client := NewImageClient(server.URL+"/api/v1/services/aigc/multimodal-generation/generation", "test-key", "wan2.7-image-pro", "", uploadDir, time.Second)
	data, err := client.Generate(context.Background(), "/uploads/selfies/me.png", "", "合成打卡照")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("generated image is not a png: %v", err)
	}
	for _, want := range []string{
		`"model":"wan2.7-image-pro"`,
		`"thinking_mode":true`,
		`"watermark":false`,
		`"text":"合成打卡照"`,
		`"image":"data:image/png;base64,`,
	} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("request body missing %s: %s", want, gotBody)
		}
	}
}

func TestDecodeFirstImageFindsNestedURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(tinyPNG(t))
	}))
	defer server.Close()

	value := map[string]any{
		"output": map[string]any{
			"choices": []any{
				map[string]any{
					"message": map[string]any{
						"content": []any{map[string]any{"image": server.URL}},
					},
				},
			},
		},
	}
	data, err := decodeFirstImage(context.Background(), value)
	if err != nil {
		t.Fatalf("decodeFirstImage() error = %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("decoded image is not a png: %v", err)
	}
}

func tinyPNG(t *testing.T) []byte {
	t.Helper()
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ioReadAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}
