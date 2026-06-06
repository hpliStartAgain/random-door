package service

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeAdminStorage struct {
	saveBytes func([]byte, string, string) (string, error)
}

func (s *fakeAdminStorage) SaveBytes(data []byte, subDir, ext string) (string, error) {
	if s.saveBytes == nil {
		return "/uploads/admin_imports/test" + ext, nil
	}
	return s.saveBytes(data, subDir, ext)
}

func TestAdminServiceImportImageURLStoresLocalImage(t *testing.T) {
	imageBytes := tinyAdminPNG(t)
	var savedSubDir string
	var savedExt string
	var savedLen int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(imageBytes)
	}))
	defer server.Close()

	svc := NewAdminService(nil, &fakeAdminStorage{
		saveBytes: func(data []byte, subDir, ext string) (string, error) {
			savedSubDir = subDir
			savedExt = ext
			savedLen = len(data)
			return "/uploads/admin_imports/imported.png", nil
		},
	})
	svc.client = server.Client()

	got, err := svc.ImportImageURL(context.Background(), server.URL+"/asset")
	if err != nil {
		t.Fatalf("ImportImageURL() error = %v", err)
	}
	if got != "/uploads/admin_imports/imported.png" {
		t.Fatalf("ImportImageURL() = %q", got)
	}
	if savedSubDir != "admin_imports" || savedExt != ".png" || savedLen != len(imageBytes) {
		t.Fatalf("saved args = %q %q %d, want admin_imports .png %d", savedSubDir, savedExt, savedLen, len(imageBytes))
	}
}

func TestAdminServiceImportImageURLRejectsFakeImageBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte("not actually an image"))
	}))
	defer server.Close()

	svc := NewAdminService(nil, &fakeAdminStorage{})
	svc.client = server.Client()

	_, err := svc.ImportImageURL(context.Background(), server.URL+"/fake.jpg")
	if !errors.Is(err, ErrInvalidParam) || ClientMessage(err) != "remote image type not supported" {
		t.Fatalf("ImportImageURL() error = %v, want invalid remote image type", err)
	}
}

func tinyAdminPNG(t *testing.T) []byte {
	t.Helper()
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}
	return data
}
