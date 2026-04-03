package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveBaseImage_MoonbitFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "moonbit": {
    "image": "ubuntu:24.04",
    "install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	image, install, err := resolveBaseImage("moonbit", "", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image != "ubuntu:24.04" {
		t.Fatalf("unexpected image: %s", image)
	}
	if install != "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash" {
		t.Fatalf("unexpected install command: %s", install)
	}
}

func TestResolveBaseImage_ExplicitBaseImageOverridesConfig(t *testing.T) {
	image, install, err := resolveBaseImage("moonbit", "ubuntu:latest", "codespacegen.base-images.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image != "ubuntu:latest" {
		t.Fatalf("unexpected image: %s", image)
	}
	if install != "" {
		t.Fatalf("expected empty install command, got: %s", install)
	}
}

func TestResolveBaseImage_UnsupportedLanguageReturnsError(t *testing.T) {
	_, _, err := resolveBaseImage("unknownlang", "", "codespacegen.base-images.json")
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
