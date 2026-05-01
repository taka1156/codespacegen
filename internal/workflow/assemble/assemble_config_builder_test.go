package assemble

import (
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/utils"
)

func TestResolveBaseImage_UsesLanguageEntry(t *testing.T) {
	entries := map[string]entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := resolveBaseImage("python", entries, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "python:3.12" {
		t.Errorf("Image: got %q, want %q", got.Image, "python:3.12")
	}
}

func TestResolveBaseImage_LanguageKeyIsCaseInsensitive(t *testing.T) {
	entries := map[string]entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := resolveBaseImage("Python", entries, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "python:3.12" {
		t.Errorf("Image: got %q, want %q", got.Image, "python:3.12")
	}
}

func TestResolveBaseImage_FallsBackToDefaultImageWhenLanguageEmpty(t *testing.T) {
	got, err := resolveBaseImage("", nil, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "alpine:latest" {
		t.Errorf("Image: got %q, want %q", got.Image, "alpine:latest")
	}
}

func TestResolveBaseImage_FallsBackToEntityDefaultImageWhenBothEmpty(t *testing.T) {
	got, err := resolveBaseImage("", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != entity.DefaultImage {
		t.Errorf("Image: got %q, want %q", got.Image, entity.DefaultImage)
	}
}

func TestResolveBaseImage_ErrorWhenLanguageNotInEntries(t *testing.T) {
	_, err := resolveBaseImage("rust", map[string]entity.LangEntry{}, "alpine:latest")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveBaseImage_ErrorWhenEntryImageIsEmpty(t *testing.T) {
	entries := map[string]entity.LangEntry{
		"rust": {Image: ""},
	}
	_, err := resolveBaseImage("rust", entries, "alpine:latest")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveTimezone_PromptTakesPriority(t *testing.T) {
	got := resolveTimezone("Asia/Tokyo", "America/New_York", utils.Ptr("Europe/London"), "UTC")
	if got != "Asia/Tokyo" {
		t.Errorf("got %q, want %q", got, "Asia/Tokyo")
	}
}

func TestResolveTimezone_ExplicitUsedWhenPromptEmpty(t *testing.T) {
	got := resolveTimezone("", "America/New_York", utils.Ptr("Europe/London"), "UTC")
	if got != "America/New_York" {
		t.Errorf("got %q, want %q", got, "America/New_York")
	}
}

func TestResolveTimezone_ConfigUsedWhenPromptAndExplicitEmpty(t *testing.T) {
	got := resolveTimezone("", "", utils.Ptr("Europe/London"), "UTC")
	if got != "Europe/London" {
		t.Errorf("got %q, want %q", got, "Europe/London")
	}
}

func TestResolveTimezone_DefaultUsedWhenAllEmpty(t *testing.T) {
	got := resolveTimezone("", "", nil, "UTC")
	if got != "UTC" {
		t.Errorf("got %q, want %q", got, "UTC")
	}
}

func TestResolveTimezone_EntityDefaultUsedWhenEverythingEmpty(t *testing.T) {
	got := resolveTimezone("", "", nil, "")
	if got != entity.DefaultTimezone {
		t.Errorf("got %q, want %q", got, entity.DefaultTimezone)
	}
}

func TestMergeOsModules_ReturnsBaseWhenLinuxPackagesNil(t *testing.T) {
	base := entity.OsModules{
		AlpineModules:     []string{"bash"},
		DebianLikeModules: []string{"bash"},
	}
	got := mergeOsModules(base, nil)
	if len(got.AlpineModules) != 1 || got.AlpineModules[0] != "bash" {
		t.Errorf("AlpineModules: got %v, want [bash]", got.AlpineModules)
	}
}

func TestMergeOsModules_AppendsLinuxPackages(t *testing.T) {
	base := entity.OsModules{
		AlpineModules:     []string{"bash"},
		DebianLikeModules: []string{"bash"},
	}
	pkgs := []entity.LinuxPackage{"git", "curl"}
	got := mergeOsModules(base, &pkgs)
	if len(got.AlpineModules) != 3 {
		t.Errorf("AlpineModules len: got %d, want 3", len(got.AlpineModules))
	}
	if len(got.DebianLikeModules) != 3 {
		t.Errorf("DebianLikeModules len: got %d, want 3", len(got.DebianLikeModules))
	}
}

func TestMergeOsModules_DeduplicatesPackages(t *testing.T) {
	base := entity.OsModules{
		AlpineModules:     []string{"bash", "git"},
		DebianLikeModules: []string{"bash", "git"},
	}
	pkgs := []entity.LinuxPackage{"git", "curl"}
	got := mergeOsModules(base, &pkgs)
	if len(got.AlpineModules) != 3 {
		t.Errorf("AlpineModules len: got %d, want 3 (bash, git, curl)", len(got.AlpineModules))
	}
	if len(got.DebianLikeModules) != 3 {
		t.Errorf("DebianLikeModules len: got %d, want 3 (bash, git, curl)", len(got.DebianLikeModules))
	}
}

func TestBuildCodespaceConfig_NilRunCommandAndExtensionsDoNotPanic(t *testing.T) {
	acc := NewAssembleCodespaceConfig(&fakeCodespacePromptResolver{
		projectName:     "p",
		serviceName:     "s",
		workspaceFolder: "/w",
	})

	setting := entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}
	clientConfig := entity.ClientConfig{
		ComposeFile: utils.Ptr("docker-compose.yaml"),
	}
	jsonConfig := entity.JsonConfig{Common: &entity.CommonEntry{}}

	got, err := acc.Resolve(clientConfig, setting, jsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.RunCommand != "" {
		t.Errorf("RunCommand: got %q, want empty", got.RunCommand)
	}
	if got.VSCodeExtensions == nil {
		t.Error("VSCodeExtensions should not be nil")
	}
	if len(got.VSCodeExtensions) != 0 {
		t.Errorf("VSCodeExtensions: got %v, want []", got.VSCodeExtensions)
	}
}

func TestBuildCodespaceConfig_NoLanguageUsesDefaultImage(t *testing.T) {
	acc := NewAssembleCodespaceConfig(&fakeCodespacePromptResolver{
		projectName:     "p",
		serviceName:     "s",
		workspaceFolder: "/w",
		language:        "",
	})

	setting := entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}
	clientConfig := entity.ClientConfig{
		ComposeFile: utils.Ptr("docker-compose.yaml"),
	}
	jsonConfig := entity.JsonConfig{Common: &entity.CommonEntry{}}

	got, err := acc.Resolve(clientConfig, setting, jsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.BaseImage != "alpine:latest" {
		t.Errorf("BaseImage: got %q, want %q", got.BaseImage, "alpine:latest")
	}
}

func TestResolvePort_PromptTakesPriority(t *testing.T) {
	got := resolvePort("3000:3000", "8080:8080")
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestResolvePort_ExplicitUsedWhenPromptEmpty(t *testing.T) {
	got := resolvePort("", "8080:8080")
	if got != "8080:8080" {
		t.Errorf("got %q, want %q", got, "8080:8080")
	}
}

func TestResolvePort_EmptyWhenBothEmpty(t *testing.T) {
	got := resolvePort("", "")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
