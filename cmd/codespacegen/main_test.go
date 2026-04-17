package main

import (
	"os"
	"path/filepath"
	"testing"

	"codespacegen/internal/domain/entity"
)

func TestResolveTimezone_UsesDefaultWhenEmpty(t *testing.T) {
	timezone, err := resolveTimezone("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if timezone != entity.DefaultTimezone {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveTimezone_TrimsExplicitValue(t *testing.T) {
	timezone, err := resolveTimezone("  Europe/Berlin  ", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if timezone != "Europe/Berlin" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveTimezone_UsesConfigTimezoneWhenFlagEmpty(t *testing.T) {
	timezone, err := resolveTimezone("", "  America/Los_Angeles  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if timezone != "America/Los_Angeles" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveBaseImage_MoonbitFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "moonbit": {
    "image": "ubuntu:24.04",
		"install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash",
		"timezone": "Europe/Berlin"
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	image, install, timezone, err := resolveBaseImage("moonbit", "", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image != "ubuntu:24.04" {
		t.Fatalf("unexpected image: %s", image)
	}
	if install != "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash" {
		t.Fatalf("unexpected install command: %s", install)
	}
	if timezone != "Europe/Berlin" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveBaseImage_ExplicitBaseImageOverridesConfig(t *testing.T) {
	image, install, timezone, err := resolveBaseImage("moonbit", "ubuntu:latest", "codespacegen.base-images.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image != "ubuntu:latest" {
		t.Fatalf("unexpected image: %s", image)
	}
	if install != "" {
		t.Fatalf("expected empty install command, got: %s", install)
	}
	if timezone != "" {
		t.Fatalf("expected empty timezone, got: %s", timezone)
	}
}

func TestResolveBaseImage_UsesDefaultImageWhenOnlyTimezoneIsOverridden(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "go": {
    "timezone": "UTC"
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	image, install, timezone, err := resolveBaseImage("go", "", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image != "golang:1.24-alpine" {
		t.Fatalf("unexpected image: %s", image)
	}
	if install != "" {
		t.Fatalf("unexpected install command: %s", install)
	}
	if timezone != "UTC" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveBaseImage_UnsupportedLanguageReturnsError(t *testing.T) {
	_, _, _, err := resolveBaseImage("unknownlang", "", "codespacegen.base-images.json")
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
