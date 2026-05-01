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
