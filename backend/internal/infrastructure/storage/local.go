package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage persists objects on the local filesystem under a configurable
// base directory. It is safe for concurrent use at the OS level (file
// operations are atomic per call) but does not add internal locking.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a LocalStorage rooted at basePath. The directory (and
// any missing parents) are created with 0o750 permissions if they do not exist.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("storage: base path is required")
	}
	if err := os.MkdirAll(basePath, 0o750); err != nil {
		return nil, fmt.Errorf("storage: create base dir %q: %w", basePath, err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

// fullPath joins and cleans the key against the base path, refusing any key
// that escapes the base directory (defence against path traversal).
func (s *LocalStorage) fullPath(key string) (string, error) {
	cleaned := filepath.Clean(key)
	if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("storage: invalid key %q", key)
	}
	return filepath.Join(s.basePath, cleaned), nil
}

// Upload writes the reader's contents to disk under key.
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, _ string) error {
	path, err := s.fullPath(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("storage: create dir for %q: %w", key, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("storage: create %q: %w", key, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("storage: write %q: %w", key, err)
	}
	return nil
}

// Download returns a ReadCloser for the object at key.
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	path, err := s.fullPath(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("storage: not found %q: %w", key, err)
		}
		return nil, fmt.Errorf("storage: open %q: %w", key, err)
	}
	return f, nil
}

// Delete removes the object at key. Missing files are treated as success.
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	path, err := s.fullPath(key)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: delete %q: %w", key, err)
	}
	return nil
}
