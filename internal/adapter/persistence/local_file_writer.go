package persistence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type LocalFileWriter struct {
}

func NewLocalFileWriter() *LocalFileWriter {
	return &LocalFileWriter{}
}

func (w *LocalFileWriter) Write(relativePath string, content string, overwrite bool) error {
	fullPath := relativePath

	if !overwrite {
		if _, err := os.Stat(fullPath); err == nil {
			return fmt.Errorf("file already exists: %s", fullPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0o644)
}
