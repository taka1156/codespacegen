package generator

import (
	"strings"
	"testing"

	"codespacegen/internal/domain/entity"
)

func TestDefaultTemplateGenerator_Generate_UsesApkForAlpineImage(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "alpine:latest",
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "apk add --no-cache") {
		t.Fatal("expected Dockerfile to use apk for alpine image")
	}
	if !strings.Contains(dockerfile, "ca-certificates") {
		t.Fatal("expected Dockerfile to install ca-certificates for alpine image")
	}
	if strings.Contains(dockerfile, "apt-get install") {
		t.Fatal("expected Dockerfile not to use apt-get for alpine image")
	}
}

func TestDefaultTemplateGenerator_Generate_UsesAptForUbuntuImage(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "apt-get install -y --no-install-recommends") {
		t.Fatal("expected Dockerfile to use apt-get for ubuntu image")
	}
	if !strings.Contains(dockerfile, "ca-certificates") {
		t.Fatal("expected Dockerfile to install ca-certificates for ubuntu image")
	}
	if strings.Contains(dockerfile, "apk add --no-cache") {
		t.Fatal("expected Dockerfile not to use apk for ubuntu image")
	}
}

func TestDefaultTemplateGenerator_Generate_EmbedsInstallCommand(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		ComposeFileName: "docker-compose.yaml",
		InstallCommand:  "apt-get update && apt-get install -y build-essential",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "RUN apt-get update && apt-get install -y build-essential") {
		t.Fatal("expected Dockerfile to include install command")
	}
}

func findGeneratedFile(t *testing.T, files []entity.GeneratedFile, relativePath string) string {
	t.Helper()
	for _, f := range files {
		if f.RelativePath == relativePath {
			return f.Content
		}
	}
	t.Fatalf("generated file not found: %s", relativePath)
	return ""
}
