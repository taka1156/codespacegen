package utils

import (
	"testing"
)

func TestPtr_ReturnsPointerToValue(t *testing.T) {
	got := Ptr("hello")
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != "hello" {
		t.Errorf("got %q, want %q", *got, "hello")
	}
}

func TestPtr_IntPointer(t *testing.T) {
	got := Ptr(42)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != 42 {
		t.Errorf("got %d, want %d", *got, 42)
	}
}

func TestPtr_BoolPointer(t *testing.T) {
	got := Ptr(true)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != true {
		t.Errorf("got %v, want true", *got)
	}
}

func TestPtr_ReturnsUniquePointers(t *testing.T) {
	p1 := Ptr("value")
	p2 := Ptr("value")
	if p1 == p2 {
		t.Error("expected different pointer addresses for separate calls")
	}
}

func TestNormalizePortMapping_PortOnlyNormalizesToMapping(t *testing.T) {
	got, err := NormalizePortMapping("3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestNormalizePortMapping_FullMappingPassesThrough(t *testing.T) {
	got, err := NormalizePortMapping("8080:9090")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "8080:9090" {
		t.Errorf("got %q, want %q", got, "8080:9090")
	}
}

func TestNormalizePortMapping_LeadingTrailingSpacesTrimmed(t *testing.T) {
	got, err := NormalizePortMapping("  3000  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestNormalizePortMapping_InvalidValueReturnsError(t *testing.T) {
	_, err := NormalizePortMapping("bad")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizePortMapping_EmptyValueReturnsError(t *testing.T) {
	_, err := NormalizePortMapping("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizePortMapping_AlphaNumericReturnsError(t *testing.T) {
	_, err := NormalizePortMapping("3000abc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
