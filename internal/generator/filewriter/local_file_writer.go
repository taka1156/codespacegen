package filewriter

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

func (w *LocalFileWriter) Write(path string, content string, overwrite bool) error {
	fullPath := path

	if !overwrite {
		if _, err := os.Stat(fullPath); err == nil {
			return fmt.Errorf("file already exists: %s", fullPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(fullPath), "."+filepath.Base(fullPath)+".tmp-*")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Chmod(0o644); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, fullPath); err != nil {
		return err
	}

	return nil
}
