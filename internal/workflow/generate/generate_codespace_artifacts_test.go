package generate

import (
	"errors"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

type fakeGenerator struct {
	files []entity.GeneratedFile
	err   error
}

func (g fakeGenerator) Generate(_ entity.CodespaceConfig) ([]entity.GeneratedFile, error) {
	return g.files, g.err
}

type writeCall struct {
	path      string
	content   string
	overwrite bool
}

type fakeWriter struct {
	calls []writeCall
	err   error
}

func (w *fakeWriter) Write(path string, content string, overwrite bool) error {
	if w.err != nil {
		return w.err
	}
	w.calls = append(w.calls, writeCall{path: path, content: content, overwrite: overwrite})
	return nil
}

func TestGenerateCodespaceArtifacts_Execute_WritesAllFiles(t *testing.T) {
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "debian:bookworm",
		ComposeFileName: "docker-compose.yaml",
	}

	gen := fakeGenerator{files: []entity.GeneratedFile{
		{RelativePath: "Dockerfile", Content: "A"},
		{RelativePath: "devcontainer.json", Content: "B"},
	}}
	writer := &fakeWriter{}
	uc := NewGenerateCodespaceArtifacts(gen, writer)

	if err := uc.Execute(cfg, true, ".devcontainer"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(writer.calls) != 2 {
		t.Fatalf("expected 2 writes, got %d", len(writer.calls))
	}
	if writer.calls[0].path != ".devcontainer/Dockerfile" {
		t.Fatalf("unexpected first path: %s", writer.calls[0].path)
	}
	if !writer.calls[0].overwrite {
		t.Fatal("expected overwrite to be true")
	}
}

func TestGenerateCodespaceArtifacts_Execute_ReturnsWriteError(t *testing.T) {
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "debian:bookworm",
		ComposeFileName: "docker-compose.yaml",
	}

	gen := fakeGenerator{files: []entity.GeneratedFile{{RelativePath: "Dockerfile", Content: "A"}}}
	writer := &fakeWriter{err: errors.New("write failed")}
	uc := NewGenerateCodespaceArtifacts(gen, writer)

	err := uc.Execute(cfg, false, ".devcontainer")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
