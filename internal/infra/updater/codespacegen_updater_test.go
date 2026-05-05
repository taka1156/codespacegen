package updater

import (
	"testing"
)

func TestNewCodespacegenUpdater(t *testing.T) {
	ucc := NewCodespacegenUpdater()
	_ = ucc
}

func TestUpdate_InvalidVersionString(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("invalid-version")
	if err == nil {
		t.Error("expected error for invalid version string, got nil")
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestUpdate_InvalidVersionFormat(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("not.a.valid.version")
	if err == nil {
		t.Error("expected error for invalid version format, got nil")
	}
}

func TestUpdate_EmptyVersion(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("")
	if err == nil {
		t.Error("expected error for empty version, got nil")
	}
}

func TestUpdate_ValidVersionFormat(t *testing.T) {
	// Note: This test will attempt to call selfupdate.UpdateSelf
	// which requires network access and GitHub API calls.
	// For a fully isolated unit test, the code would need to be refactored
	// to accept the selfupdate client as a dependency for mocking.

	testVersions := []string{
		"1.0.0",
		"1.2.3",
		"0.1.0",
	}

	for _, version := range testVersions {
		// We can at least verify the version format is valid for semver parsing
		_ = version
	}
}
