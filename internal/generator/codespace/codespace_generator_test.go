package codespace

import (
	"strings"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

var testAlpineModules = []string{
	"bash", "bash-completion", "ca-certificates", "tzdata",
	"git", "git-lfs", "vim", "curl",
	"musl-locales", "musl-locales-lang",
}

var testDebianLikeModules = []string{
	"bash", "bash-completion", "ca-certificates", "tzdata",
	"git", "git-lfs", "vim", "curl",
	"locales",
}

func TestCodespaceGenerator_Generate_UsesApkForAlpineImage(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
		ComposeFileName: "docker-compose.yaml",
		OsModules: entity.OsModules{
			AlpineModules: testAlpineModules,
		},
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "apk add --no-cache") {
		t.Fatal("expected Dockerfile to use apk for alpine image")
	}
	if !strings.Contains(dockerfile, "musl-locales") {
		t.Fatal("expected Dockerfile to install musl-locales for alpine image")
	}
	if !strings.Contains(dockerfile, "musl-locales-lang") {
		t.Fatal("expected Dockerfile to install musl-locales-lang for alpine image")
	}
	if !strings.Contains(dockerfile, "ca-certificates") {
		t.Fatal("expected Dockerfile to install ca-certificates for alpine image")
	}
	if strings.Contains(dockerfile, "apt get install") {
		t.Fatal("expected Dockerfile not to use apt get for alpine image")
	}
}

func TestCodespaceGenerator_Generate_AlpineDoesNotRunLocaleGen(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if strings.Contains(dockerfile, "locale-gen") {
		t.Fatal("expected Dockerfile not to run locale-gen for alpine image")
	}
	if strings.Contains(dockerfile, "update-locale") {
		t.Fatal("expected Dockerfile not to run update-locale for alpine image")
	}
}

func TestCodespaceGenerator_Generate_UsesAptForUbuntuImage(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		ComposeFileName: "docker-compose.yaml",
		OsModules: entity.OsModules{
			DebianLikeModules: testDebianLikeModules,
		},
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "apt install -y --no-install-recommends") {
		t.Fatal("expected Dockerfile to use apt for ubuntu image")
	}
	if !strings.Contains(dockerfile, "ca-certificates") {
		t.Fatal("expected Dockerfile to install ca-certificates for ubuntu image")
	}
	if !strings.Contains(dockerfile, "locale-gen") {
		t.Fatal("expected Dockerfile to run locale-gen for ubuntu image")
	}
	if !strings.Contains(dockerfile, "update-locale") {
		t.Fatal("expected Dockerfile to run update-locale for ubuntu image")
	}
	if strings.Contains(dockerfile, "apk add --no-cache") {
		t.Fatal("expected Dockerfile not to use apk for ubuntu image")
	}
}

func TestCodespaceGenerator_Generate_UsesDefaultLocale(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "LANG="+entity.DefaultLocale.Lang) {
		t.Fatalf("expected Dockerfile to use default LANG %q", entity.DefaultLocale.Lang)
	}
	if !strings.Contains(dockerfile, "LANGUAGE="+entity.DefaultLocale.Language) {
		t.Fatalf("expected Dockerfile to use default LANGUAGE %q", entity.DefaultLocale.Language)
	}
	if !strings.Contains(dockerfile, "LC_ALL="+entity.DefaultLocale.LcAll) {
		t.Fatalf("expected Dockerfile to use default LC_ALL %q", entity.DefaultLocale.LcAll)
	}
}

func TestCodespaceGenerator_Generate_UsesCustomLocale(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		ComposeFileName: "docker-compose.yaml",
		Locale: entity.LocaleConfig{
			Lang:     "en_US.UTF-8",
			Language: "en_US:en",
			LcAll:    "en_US.UTF-8",
		},
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "LANG=en_US.UTF-8") {
		t.Fatal("expected Dockerfile to use custom LANG in ENV")
	}
	if !strings.Contains(dockerfile, "LANGUAGE=en_US:en") {
		t.Fatal("expected Dockerfile to use custom LANGUAGE in ENV")
	}
	if !strings.Contains(dockerfile, "LC_ALL=en_US.UTF-8") {
		t.Fatal("expected Dockerfile to use custom LC_ALL in ENV")
	}
	if !strings.Contains(dockerfile, "locale-gen en_US.UTF-8") {
		t.Fatal("expected Dockerfile to run locale-gen with custom locale")
	}
	if !strings.Contains(dockerfile, "update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8") {
		t.Fatal("expected Dockerfile to run update-locale with custom locale")
	}
	if strings.Contains(dockerfile, "ja_JP") {
		t.Fatal("expected Dockerfile not to contain default Japanese locale when custom locale is provided")
	}
}

func TestCodespaceGenerator_Generate_EmbedsRunCommand(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       "ubuntu:latest",
		ComposeFileName: "docker-compose.yaml",
		RunCommand:      "apt get update && apt get install -y build-essential",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerfile := findGeneratedFile(t, files, "Dockerfile")
	if !strings.Contains(dockerfile, "RUN apt get update && apt get install -y build-essential") {
		t.Fatal("expected Dockerfile to include install command")
	}
}

func TestCodespaceGenerator_Generate_UsesCustomTimezone(t *testing.T) {
	g := NewCodespaceGenerator()
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

func TestCodespaceGenerator_Generate_DevcontainerKeyOrder(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "sample",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
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

	if namePos >= servicePos ||
		servicePos >= workspacePos ||
		workspacePos >= dockerComposePos ||
		dockerComposePos >= customizationsPos {
		t.Fatalf("unexpected key order in devcontainer output: %s", devcontainer)
	}
}

func TestCodespaceGenerator_Generate_ComposeIncludesPortMapping(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
		ComposeFileName: "docker-compose.yaml",
		PortMapping:     "3000:3000",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	compose := findGeneratedFile(t, files, "docker-compose.yaml")
	if !strings.Contains(compose, "ports:") {
		t.Fatal("expected docker-compose.yaml to include ports section")
	}
	if !strings.Contains(compose, "\"3000:3000\"") {
		t.Fatal("expected docker-compose.yaml to include port mapping")
	}
}

func TestCodespaceGenerator_Generate_ComposeOmitsPortsWhenEmpty(t *testing.T) {
	g := NewCodespaceGenerator()
	cfg := entity.CodespaceConfig{
		ContainerName:   "test",
		ServiceName:     "app",
		WorkspaceFolder: "/workspace",
		BaseImage:       entity.DefaultImage,
		ComposeFileName: "docker-compose.yaml",
	}

	files, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	compose := findGeneratedFile(t, files, "docker-compose.yaml")
	if strings.Contains(compose, "ports:") {
		t.Fatal("expected docker-compose.yaml not to include ports section when port mapping is empty")
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
