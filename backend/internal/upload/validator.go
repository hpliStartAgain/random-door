package upload

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var (
	ErrFileTooLarge    = errors.New("file exceeds maximum size")
	ErrUnsupportedType = errors.New("unsupported file type")
)

type Validator struct {
	MaxSizeBytes int64
	AllowedTypes map[string]bool
}

func NewValidator(maxSizeMB int, allowedTypes []string) *Validator {
	types := make(map[string]bool)
	for _, t := range allowedTypes {
		types[strings.ToLower(strings.TrimSpace(t))] = true
	}
	return &Validator{
		MaxSizeBytes: int64(maxSizeMB) * 1024 * 1024,
		AllowedTypes: types,
	}
}

func (v *Validator) Validate(fh *multipart.FileHeader) error {
	if fh.Size > v.MaxSizeBytes {
		return ErrFileTooLarge
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))
	if !v.AllowedTypes[ext] {
		return ErrUnsupportedType
	}
	return nil
}
