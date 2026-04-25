package entity

import (
	"errors"
	"testing"
)

func TestErrInvalidConfig_ReturnsErrorWithMessage(t *testing.T) {
	err := ErrInvalidConfig("container name is required")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Error() != "container name is required" {
		t.Errorf("got %q, want %q", err.Error(), "container name is required")
	}
}

func TestErrInvalidConfig_IsNotWrapped(t *testing.T) {
	err := ErrInvalidConfig("some error")
	var target invalidConfigError
	if !errors.As(err, &target) {
		t.Error("expected error to be of type invalidConfigError")
	}
}
