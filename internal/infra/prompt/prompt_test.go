package prompt

import (
	"strings"
	"testing"
)

func newResolver(input string) *CodespacegenPrompter {
	return NewCodespacegenPrompter(strings.NewReader(input))
}

func TestPromptLanguage_UsesExplicitValueWhenUserAccepts(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptLanguage("Python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "python" {
		t.Errorf("got %q, want %q", got, "python")
	}
}

func TestPromptLanguage_UserOverridesWithInput(t *testing.T) {
	r := newResolver("Rust\n")
	got, err := r.PromptLanguage("python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "rust" {
		t.Errorf("got %q, want %q", got, "rust")
	}
}

func TestPromptLanguage_NoExplicitUserInputsValue(t *testing.T) {
	r := newResolver("go\n")
	got, err := r.PromptLanguage("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "go" {
		t.Errorf("got %q, want %q", got, "go")
	}
}

func TestPromptLanguage_NoExplicitUserAcceptsEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptLanguage("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestPromptWorkspaceFolder_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptWorkspaceFolder("/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/app" {
		t.Errorf("got %q, want %q", got, "/app")
	}
}

func TestPromptWorkspaceFolder_DefaultsToWorkspaceWhenEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptWorkspaceFolder("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/workspace" {
		t.Errorf("got %q, want %q", got, "/workspace")
	}
}

func TestPromptWorkspaceFolder_UserOverrides(t *testing.T) {
	r := newResolver("/custom\n")
	got, err := r.PromptWorkspaceFolder("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom" {
		t.Errorf("got %q, want %q", got, "/custom")
	}
}

func TestPromptServiceName_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptServiceName("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "myapp" {
		t.Errorf("got %q, want %q", got, "myapp")
	}
}

func TestPromptServiceName_DefaultsToAppWhenEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptServiceName("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "app" {
		t.Errorf("got %q, want %q", got, "app")
	}
}

func TestPromptServiceName_UserOverrides(t *testing.T) {
	r := newResolver("backend\n")
	got, err := r.PromptServiceName("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "backend" {
		t.Errorf("got %q, want %q", got, "backend")
	}
}

func TestPromptTimezone_UsesDefaultWhenUserAccepts(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptTimezone("UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "UTC" {
		t.Errorf("got %q, want %q", got, "UTC")
	}
}

func TestPromptTimezone_UserOverridesDefault(t *testing.T) {
	r := newResolver("Asia/Tokyo\n")
	got, err := r.PromptTimezone("UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Asia/Tokyo" {
		t.Errorf("got %q, want %q", got, "Asia/Tokyo")
	}
}

func TestPromptTimezone_EmptyDefaultAndEmptyInputReturnsEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptTimezone("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestPromptProjectName_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptProjectName("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "myproject" {
		t.Errorf("got %q, want %q", got, "myproject")
	}
}

func TestPromptProjectName_UserInputsValue(t *testing.T) {
	r := newResolver("hello\n")
	got, err := r.PromptProjectName("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestPromptProjectName_EOFWithNoInputReturnsError(t *testing.T) {
	r := newResolver("") // EOF immediately
	_, err := r.PromptProjectName("")
	if err == nil {
		t.Fatal("expected error for empty project name at EOF, got nil")
	}
}

func TestPromptPortMapping_EmptyPorts(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptPortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestPromptPortMapping_NormalizesPortOnlyToMapping(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptPortMapping("3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestPromptPortMapping_AcceptsFullMappingFormat(t *testing.T) {
	r := newResolver("\n")
	got, err := r.PromptPortMapping("8080:9090")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "8080:9090" {
		t.Errorf("got %q, want %q", got, "8080:9090")
	}
}

func TestPromptPortMapping_UserInputsPort(t *testing.T) {
	r := newResolver("4000\n")
	got, err := r.PromptPortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "4000:4000" {
		t.Errorf("got %q, want %q", got, "4000:4000")
	}
}

func TestPromptPortMapping_RetriesOnInvalidThenAcceptsValid(t *testing.T) {
	r := newResolver("bad\n3000\n")
	got, err := r.PromptPortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestPromptPortMapping_InvalidDefaultPortRetriesUntilValid(t *testing.T) {
	r := newResolver("\n5000\n")
	got, err := r.PromptPortMapping("bad")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "5000:5000" {
		t.Errorf("got %q, want %q", got, "5000:5000")
	}
}
