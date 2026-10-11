package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage хранит файлы на диске внутри root. Пути — относительные, со слешами.
type LocalStorage struct {
	root string
}

func NewLocalStorage(root string) *LocalStorage {
	return &LocalStorage{root: root}
}

func (s *LocalStorage) Save(relPath string, r io.Reader) error {
	full, err := s.resolve(relPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(full)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(full)
		return fmt.Errorf("write file: %w", err)
	}
	return f.Close()
}

func (s *LocalStorage) Open(relPath string) (io.ReadCloser, error) {
	full, err := s.resolve(relPath)
	if err != nil {
		return nil, err
	}
	return os.Open(full)
}

func (s *LocalStorage) Remove(relPath string) error {
	full, err := s.resolve(relPath)
	if err != nil {
		return err
	}
	return os.Remove(full)
}

// resolve не даёт выйти за пределы root через "../".
func (s *LocalStorage) resolve(relPath string) (string, error) {
	full := filepath.Join(s.root, filepath.FromSlash(relPath))
	root, _ := filepath.Abs(s.root)
	abs, _ := filepath.Abs(full)
	if abs != root && !strings.HasPrefix(abs, root+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes storage root", relPath)
	}
	return full, nil
}
