package filewriter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalFileWriter_Write_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w := NewLocalFileWriter()
	if err := w.Write(path, "hello", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("got %q, want %q", string(got), "hello")
	}
}

func TestLocalFileWriter_Write_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "out.txt")

	w := NewLocalFileWriter()
	if err := w.Write(path, "content", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not found: %v", err)
	}
}

func TestLocalFileWriter_Write_ErrorIfExistsWithoutForce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w := NewLocalFileWriter()
	if err := w.Write(path, "first", false); err != nil {
		t.Fatalf("unexpected error on first write: %v", err)
	}

	err := w.Write(path, "second", false)
	if err == nil {
		t.Fatal("expected error when overwrite=false, got nil")
	}
}

func TestLocalFileWriter_Write_OverwritesWithForce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w := NewLocalFileWriter()
	if err := w.Write(path, "first", false); err != nil {
		t.Fatalf("unexpected error on first write: %v", err)
	}
	if err := w.Write(path, "second", true); err != nil {
		t.Fatalf("unexpected error on overwrite: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(got) != "second" {
		t.Errorf("got %q, want %q", string(got), "second")
	}
}
