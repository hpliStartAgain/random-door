package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Storage struct {
	BaseDir string // e.g. "./uploads"
}

func NewStorage(baseDir string) *Storage {
	return &Storage{BaseDir: baseDir}
}

// Save stores a file under baseDir/subDir with a UUID-based name.
// Returns the relative URL path (e.g. "/uploads/selfies/abc.jpg").
func (s *Storage) Save(fh *multipart.FileHeader, subDir string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	fileName := uuid.New().String() + ext

	dirPath := filepath.Join(s.BaseDir, subDir)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}

	fullPath := filepath.Join(dirPath, fileName)
	// Prevent path traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(s.BaseDir)) {
		return "", fmt.Errorf("invalid path")
	}

	src, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	// Return URL path
	urlPath := "/" + filepath.ToSlash(filepath.Join(filepath.Base(s.BaseDir), subDir, fileName))
	return urlPath, nil
}

// SaveBytes stores raw bytes under baseDir/subDir with a UUID-based name.
func (s *Storage) SaveBytes(data []byte, subDir, ext string) (string, error) {
	fileName := uuid.New().String() + ext

	dirPath := filepath.Join(s.BaseDir, subDir)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}

	fullPath := filepath.Join(dirPath, fileName)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	urlPath := "/" + filepath.ToSlash(filepath.Join(filepath.Base(s.BaseDir), subDir, fileName))
	return urlPath, nil
}
