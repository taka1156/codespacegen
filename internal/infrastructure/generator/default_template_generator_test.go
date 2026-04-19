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

func TestDefaultTemplateGenerator_Generate_UsesCustomTimezone(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		Timezone:        "America/New_York",
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "TZ=America/New_York") {
		t.Fatal("expected Dockerfile to use custom timezone in ENV")
	}
	if !strings.Contains(dockerfile, "/usr/share/zoneinfo/America/New_York") {
		t.Fatal("expected Dockerfile to use custom timezone in setup block")
	}
	if strings.Contains(dockerfile, entity.DefaultTimezone) {
		t.Fatal("expected Dockerfile not to keep default timezone when custom timezone is provided")
	}
}

func TestDefaultTemplateGenerator_Generate_EmbedsVSCodeExtensions(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:    "test",
		ServiceName:      "app",
		WorkspaceFolder:  "/workspace",
		BaseImage:        "alpine:latest",
		ComposeFileName:  "docker-compose.yaml",
		VSCodeExtensions: []string{"MS-CEINTL.vscode-language-pack-ja", "golang.Go", "MS-CEINTL.vscode-language-pack-ja"},
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	devcontainer := findGeneratedFile(t, files, "devcontainer.json")
	if !strings.Contains(devcontainer, "GitHub.copilot") {
		t.Fatal("expected devcontainer.json to include default GitHub.copilot extension")
	}
	if !strings.Contains(devcontainer, "MS-CEINTL.vscode-language-pack-ja") {
		t.Fatal("expected devcontainer.json to include merged common extension")
	}
	if !strings.Contains(devcontainer, "golang.Go") {
		t.Fatal("expected devcontainer.json to include language extension")
	}
	if strings.Count(devcontainer, "MS-CEINTL.vscode-language-pack-ja") != 1 {
		t.Fatal("expected devcontainer.json to deduplicate duplicated extension")
	}
}

func TestDefaultTemplateGenerator_Generate_DevcontainerKeyOrder(t *testing.T) {
	g := NewDefaultTemplateGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "sample",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "alpine:latest",
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	devcontainer := findGeneratedFile(t, files, "devcontainer.json")

	namePos := strings.Index(devcontainer, "\"name\"")
	servicePos := strings.Index(devcontainer, "\"service\"")
	workspacePos := strings.Index(devcontainer, "\"workspaceFolder\"")
	dockerComposePos := strings.Index(devcontainer, "\"dockerComposeFile\"")
	customizationsPos := strings.Index(devcontainer, "\"customizations\"")

	if namePos == -1 || servicePos == -1 || workspacePos == -1 || dockerComposePos == -1 || customizationsPos == -1 {
		t.Fatalf("expected all keys to exist in devcontainer output: %s", devcontainer)
	}

	if !(namePos < servicePos && servicePos < workspacePos && workspacePos < dockerComposePos && dockerComposePos < customizationsPos) {
		t.Fatalf("unexpected key order in devcontainer output: %s", devcontainer)
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
