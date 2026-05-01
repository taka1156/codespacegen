package prompt

import (
	"bufio"
	"strings"
	"testing"
)

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
