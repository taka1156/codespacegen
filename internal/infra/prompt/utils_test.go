package prompt

import (
	"bufio"
	"strings"
	"testing"
)

func TestNormalizePortMapping_PortOnlyNormalizesToMapping(t *testing.T) {
	got, err := normalizePortMapping("3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestNormalizePortMapping_FullMappingPassesThrough(t *testing.T) {
	got, err := normalizePortMapping("8080:9090")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "8080:9090" {
		t.Errorf("got %q, want %q", got, "8080:9090")
	}
}

func TestNormalizePortMapping_LeadingTrailingSpacesTrimmed(t *testing.T) {
	got, err := normalizePortMapping("  3000  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestNormalizePortMapping_InvalidValueReturnsError(t *testing.T) {
	_, err := normalizePortMapping("bad")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizePortMapping_EmptyValueReturnsError(t *testing.T) {
	_, err := normalizePortMapping("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizePortMapping_AlphaNumericReturnsError(t *testing.T) {
	_, err := normalizePortMapping("3000abc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func newBufReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestPromptWithDefault_EmptyInputReturnsDefault(t *testing.T) {
	got, err := promptWithDefault(newBufReader("\n"), "", "defaultVal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "defaultVal" {
		t.Errorf("got %q, want %q", got, "defaultVal")
	}
}

func TestPromptWithDefault_UserInputOverridesDefault(t *testing.T) {
	got, err := promptWithDefault(newBufReader("userInput\n"), "", "defaultVal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "userInput" {
		t.Errorf("got %q, want %q", got, "userInput")
	}
}

func TestPromptWithDefault_EOFWithNoInputReturnsDefault(t *testing.T) {
	got, err := promptWithDefault(newBufReader(""), "", "defaultVal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "defaultVal" {
		t.Errorf("got %q, want %q", got, "defaultVal")
	}
}

func TestPromptWithDefault_EOFWithInputReturnsInput(t *testing.T) {
	got, err := promptWithDefault(newBufReader("eofInput"), "", "defaultVal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "eofInput" {
		t.Errorf("got %q, want %q", got, "eofInput")
	}
}

func TestPromptWithDefault_WhitespaceOnlyInputReturnsDefault(t *testing.T) {
	got, err := promptWithDefault(newBufReader("   \n"), "", "defaultVal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "defaultVal" {
		t.Errorf("got %q, want %q", got, "defaultVal")
	}
}
